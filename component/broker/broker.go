package broker

import "context"

// Broker 是后端服务之间的高性能通信组件抽象。
//
// 设计目标：薄、快、稳定。它只提供“通信原语”，不内置业务级限流/分级/合并等策略。
// 这些策略应由 gate/game/match 等上层服务实现。
type Broker interface {

	// ID 返回实现标识 (例如 "nats(core)")。
	ID() string

	// Pub 发布一条消息 (fire-and-forget).
	Pub(ctx context.Context, subject string, data []byte, opts ...PublishOption) error

	// Sub 订阅主题。需要队列组 (Queue Group) 时请使用 SubscribeOption.
	Sub(ctx context.Context, subject string, handler Handler, opts ...SubscribeOption) (Subscription, error)

	// Request 发送请求并等待响应 (request-reply).
	Request(ctx context.Context, subject string, data []byte, opts ...RequestOption) (*Message, error)

	// Reply 回复一条请求消息。通常用于 Subscribe handler 中处理 request 后发送响应。
	// msg.Reply 字段包含请求方指定的回复地址 (通常是自动生成的 inbox).
	Reply(ctx context.Context, msg *Message, data []byte, opts ...ReplyOption) error

	// Close 立即关闭连接.
	Close() error
}

// Handler 是订阅回调。Broker 不负责并发与背压控制；handler 内部应自行处理队列/worker.
type Handler func(ctx context.Context, msg *Message)

// Subscription 表示一次订阅.
type Subscription interface {
	Unsub() error
	Close() error
}

// Message 是 Broker 层的通用消息结构.
type Message struct {
	Subject string
	Reply   string
	Header  Header
	Data    []byte
}
