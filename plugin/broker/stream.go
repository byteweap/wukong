package broker

import (
	"context"
	"time"
)

// StreamBroker 是 Broker 的扩展接口，提供 JetStream 的可靠队列能力。
//
// 适用场景：
// - 结算/资产变更（必须可靠，不能丢消息）
// - 排行榜更新（需要持久化）
// - 邮件/社交事件（离线消息可靠投递）
// - 日志/审计（需要回放）
//
// 设计原则：作为可选扩展，不影响 Broker 的"薄接口"设计。
type StreamBroker interface {
	Broker // 继承基础 Broker 能力

	// Enqueue 将消息入队（WorkQueue 模式），保证至少一次投递。
	// 适用于：异步任务、结算、资产变更等必须可靠处理的场景。
	Enqueue(ctx context.Context, subject string, data []byte, opts ...EnqueueOption) error

	// Consume 消费可靠队列（Pull 模式，需要手动 Ack）。
	// 适用于：工作队列、需要可靠处理的异步任务。
	// 注意：handler 中需要手动调用 msg.Ack() / msg.Nak() / msg.Term()。
	Consume(ctx context.Context, subject string, handler StreamHandler, opts ...ConsumeOption) (StreamSubscription, error)
}

// StreamHandler 是 JetStream 消费回调。
// 参数 msg 需要能够调用 Ack/Nak/Term 方法（具体实现由 StreamBroker 提供）。
type StreamHandler func(ctx context.Context, msg *Message)

// StreamSubscription 是 JetStream 订阅（支持 Ack/Nak）。
type StreamSubscription interface {
	Subscription
	// AckPending 返回待确认的消息数。
	AckPending() int
}

// EnqueueOption 入队选项。
type EnqueueOption func(*EnqueueOptions)

type EnqueueOptions struct {
	Header Header
	// 可选的幂等键（用于去重）
	IdempotencyKey string
}

// WithEnqueueHeader 设置入队消息头。
func WithEnqueueHeader(h Header) EnqueueOption {
	return func(o *EnqueueOptions) {
		o.Header = h
	}
}

// ConsumeOption 消费选项。
type ConsumeOption func(*ConsumeOptions)

type ConsumeOptions struct {
	// Durable 消费者名称（用于断点续传）。
	Durable string
	// Queue 队列组（用于水平扩展的竞争消费者）。
	Queue string
	// MaxAckPending 最大待确认消息数（用于背压控制）。
	MaxAckPending int
	// AckWait 确认等待时间（超时后自动重投）。
	AckWait time.Duration
}

// WithIdempotencyKey 设置幂等键（用于去重）。
func WithIdempotencyKey(key string) EnqueueOption {
	return func(o *EnqueueOptions) {
		o.IdempotencyKey = key
	}
}

// WithDurable 设置持久化消费者名称。
func WithDurable(durable string) ConsumeOption {
	return func(o *ConsumeOptions) {
		o.Durable = durable
	}
}

// WithStreamQueue 设置队列组（用于水平扩展）。
func WithStreamQueue(queue string) ConsumeOption {
	return func(o *ConsumeOptions) {
		o.Queue = queue
	}
}

// WithMaxAckPending 设置最大待确认消息数（背压控制）。
func WithMaxAckPending(n int) ConsumeOption {
	return func(o *ConsumeOptions) {
		o.MaxAckPending = n
	}
}

// WithAckWait 设置确认等待时间。
func WithAckWait(d time.Duration) ConsumeOption {
	return func(o *ConsumeOptions) {
		o.AckWait = d
	}
}
