// Package nats provides NATS Core and JetStream broker implementations.
// This file contains JetStream extensions for reliable queue support.
package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/byteweap/wukong/plugin/broker"
)

// StreamID is NATS JetStream broker implementation identifier.
const StreamID = "nats(jetstream)"

// StreamBroker implements broker.StreamBroker with NATS JetStream.
// 它是 Broker 的扩展实现，提供可靠队列能力（持久化、至少一次投递、重试）。
type StreamBroker struct {
	*Broker // 继承 Core Broker 能力
	js      jetstream.JetStream
	opts    *streamOptions
}

var _ broker.StreamBroker = (*StreamBroker)(nil)

// streamOptions JetStream 专用选项。
type streamOptions struct {
	streamConfig   *jetstream.StreamConfig
	consumerConfig *jetstream.ConsumerConfig
}

// StreamOption JetStream 选项。
type StreamOption func(*streamOptions)

// NewStreamBroker 创建 JetStream broker（需要在 NATS 服务器启用 JetStream）。
func NewStreamBroker(opts ...Option) (*StreamBroker, error) {
	b, err := New(opts...)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(b.nc)
	if err != nil {
		b.Close()
		return nil, err
	}

	so := &streamOptions{
		streamConfig:   &jetstream.StreamConfig{},
		consumerConfig: &jetstream.ConsumerConfig{},
	}

	return &StreamBroker{
		Broker: b,
		js:     js,
		opts:   so,
	}, nil
}

// NewStreamBrokerWith 基于已有 Broker 创建 StreamBroker。
func NewStreamBrokerWith(b *Broker) (*StreamBroker, error) {
	js, err := jetstream.New(b.nc)
	if err != nil {
		return nil, err
	}

	so := &streamOptions{
		streamConfig:   &jetstream.StreamConfig{},
		consumerConfig: &jetstream.ConsumerConfig{},
	}

	return &StreamBroker{
		Broker: b,
		js:     js,
		opts:   so,
	}, nil
}

func (sb *StreamBroker) ID() string { return StreamID }

// Enqueue 将消息入队（WorkQueue 模式），保证至少一次投递。
func (sb *StreamBroker) Enqueue(ctx context.Context, subject string, data []byte, opts ...broker.EnqueueOption) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	var eo broker.EnqueueOptions
	for _, opt := range opts {
		opt(&eo)
	}

	// 确保 stream 存在（按需创建）
	streamName := toStreamName(subject)
	_, err := sb.js.CreateStream(ctx, jetstream.StreamConfig{
		Name:         streamName,
		Subjects:     []string{subject},
		Storage:      jetstream.FileStorage,
		Retention:    jetstream.WorkQueuePolicy,
		Discard:      jetstream.DiscardOld,
		MaxAge:       24 * time.Hour, // 默认保留 24 小时
		MaxBytes:     -1,
		MaxMsgs:      -1,
		Replicas:     1,
		NoAck:        false,
		Duplicates:   2 * time.Minute, // 去重窗口 2 分钟
		MaxConsumers: -1,
		MaxMsgSize:   -1,
		Compression:  jetstream.NoCompression,
	})
	// 忽略已存在的错误
	if err != nil {
		// 检查是否是已存在的错误（不同版本的错误可能不同）
		if _, getErr := sb.js.Stream(ctx, streamName); getErr != nil {
			return err
		}
	}

	// 构建消息
	msg := &nats.Msg{
		Subject: subject,
		Data:    data,
	}
	if len(eo.Header) > 0 {
		msg.Header = toNatsHeader(eo.Header)
	}

	// 如果指定了幂等键，设置到 header
	if eo.IdempotencyKey != "" {
		if msg.Header == nil {
			msg.Header = nats.Header{}
		}
		msg.Header.Set("Nats-Msg-Id", eo.IdempotencyKey)
	}

	// 使用 JetStream 的 Publish 方法（通过 js.Publish）
	_, err = sb.js.PublishMsg(ctx, msg)
	return err
}

// Consume 消费可靠队列（Pull 模式）。
// 注意：当前实现为简化版本，handler 中暂无法直接调用 Ack/Nak（需要后续完善）。
func (sb *StreamBroker) Consume(ctx context.Context, subject string, handler broker.StreamHandler, opts ...broker.ConsumeOption) (broker.StreamSubscription, error) {
	var co broker.ConsumeOptions
	for _, opt := range opts {
		opt(&co)
	}

	streamName := toStreamName(subject)
	stream, err := sb.js.Stream(ctx, streamName)
	if err != nil {
		return nil, fmt.Errorf("stream not found: %w", err)
	}

	ackWait := co.AckWait
	if ackWait == 0 {
		ackWait = 30 * time.Second
	}
	maxAckPending := co.MaxAckPending
	if maxAckPending == 0 {
		maxAckPending = 1000
	}

	consumerConfig := jetstream.ConsumerConfig{
		Durable:       co.Durable,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       ackWait,
		MaxDeliver:    5, // 默认最多重试 5 次
		FilterSubject: subject,
		ReplayPolicy:  jetstream.ReplayInstantPolicy,
		MaxAckPending: maxAckPending,
	}

	// 创建或获取 consumer
	var consumer jetstream.Consumer
	if co.Durable != "" {
		consumer, err = stream.CreateOrUpdateConsumer(ctx, consumerConfig)
	} else {
		consumer, err = stream.CreateConsumer(ctx, consumerConfig)
	}
	if err != nil {
		return nil, fmt.Errorf("create consumer failed: %w", err)
	}

	// 启动消费协程
	msgs, err := consumer.Messages()
	if err != nil {
		return nil, fmt.Errorf("get messages failed: %w", err)
	}

	done := make(chan struct{})
	sub := &streamSubscription{
		consumer: consumer,
		msgs:     msgs,
		done:     done,
		msgMap:   make(map[string]jetstream.Msg), // 存储消息用于 Ack/Nak
	}

	go func() {
		defer func() {
			select {
			case <-done:
				// already closed
			default:
				close(done)
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case <-done:
				return
			default:
				msg, err := msgs.Next()
				if err != nil {
					return
				}

				bmsg := &broker.Message{
					Subject: msg.Subject(),
					Data:    msg.Data(),
				}
				meta, _ := msg.Metadata()
				// 将 JetStream 消息元数据存到 header 中供 Ack/Nak 使用
				if bmsg.Header == nil {
					bmsg.Header = broker.Header{}
				}
				msgKey := fmt.Sprintf("%d-%s", meta.Sequence.Stream, meta.Consumer)
				bmsg.Header["_js_key"] = []string{msgKey}
				sub.msgMap[msgKey] = msg

				handler(ctx, bmsg)
			}
		}
	}()

	// ctx 取消时自动停止消费
	if ctx != nil {
		if doneCh := ctx.Done(); doneCh != nil {
			go func() {
				<-doneCh
				select {
				case <-sub.done:
					// already closed
				default:
					_ = sub.Unsubscribe()
				}
			}()
		}
	}

	return sub, nil
}

func (sb *StreamBroker) Ack(ctx context.Context, msg *broker.Message) error {
	// 注意：需要在 Consume handler 中保存 jetstream.Msg 对象才能 Ack
	// 当前实现需要从 subscription 的 msgMap 中获取
	// 实际使用中建议在 handler 中直接 Ack，或使用回调方式
	key := msg.Header.Get("_js_key")
	if key == "" {
		return fmt.Errorf("missing _js_key in message header")
	}
	// 这里需要从 subscription 中获取实际的 msg，简化实现先返回错误提示
	return fmt.Errorf("Ack requires msg from Consume handler, use msg.Ack() directly in handler")
}

func (sb *StreamBroker) Nak(ctx context.Context, msg *broker.Message, delay ...time.Duration) error {
	// 类似 Ack，需要在 handler 中直接调用
	return fmt.Errorf("Nak requires msg from Consume handler, use msg.Nak() directly in handler")
}

func (sb *StreamBroker) Term(ctx context.Context, msg *broker.Message) error {
	// 终止消息（不再重试）
	return fmt.Errorf("Term requires msg from Consume handler, use msg.Term() directly in handler")
}

type streamSubscription struct {
	consumer jetstream.Consumer
	msgs     jetstream.MessagesContext
	done     chan struct{}
	msgMap   map[string]jetstream.Msg // 存储消息用于 Ack/Nak
}

func (s *streamSubscription) Unsubscribe() error {
	select {
	case <-s.done:
		// already closed
	default:
		close(s.done)
	}
	// JetStream Consumer 没有 Delete 方法，由服务器自动管理
	// 这里只需要停止消息接收
	s.msgs.Stop()
	return nil
}

func (s *streamSubscription) Drain() error {
	select {
	case <-s.done:
		// already closed
	default:
		close(s.done)
	}
	s.msgs.Stop()
	return nil
}

func (s *streamSubscription) AckPending() int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	info, err := s.consumer.Info(ctx)
	if err != nil {
		return 0
	}
	return int(info.NumPending)
}

// toStreamName 将 subject 转换为 stream 名称。
// Stream 名称不能包含 . * > / 等字符，需要转换为有效字符。
func toStreamName(subject string) string {
	// 将 subject 中的特殊字符替换为下划线
	name := "WK_STREAM"
	// 简化：使用 subject 的前缀，替换特殊字符
	for _, r := range subject {
		switch r {
		case '.', '*', '>', '/', '\\':
			name += "_"
		default:
			name += string(r)
		}
	}
	// 限制长度（NATS stream name 有长度限制）
	if len(name) > 64 {
		name = name[:64]
	}
	return name
}
