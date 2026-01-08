package broker

// PublishOption 发布选项（预留）。
type PublishOption func(*PublishOptions)

type PublishOptions struct {
	Header Header
}

// SubscribeOption 订阅选项。
type SubscribeOption func(*SubscribeOptions)

type SubscribeOptions struct {
	Queue string // Queue Group（为空表示普通订阅/广播）
}

// RequestOption 请求选项。
type RequestOption func(*RequestOptions)

type RequestOptions struct {
	Header Header
}

// WithHeader 设置消息头（Publish/Request）。
func WithHeader(h Header) PublishOption {
	return func(o *PublishOptions) {
		o.Header = h
	}
}

// WithRequestHeader 设置消息头（Request）。
func WithRequestHeader(h Header) RequestOption {
	return func(o *RequestOptions) {
		o.Header = h
	}
}

// WithQueue 设置订阅的 Queue Group（用于水平扩展的竞争消费者）。
func WithQueue(queue string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = queue
	}
}

// ReplyOption 回复选项。
type ReplyOption func(*ReplyOptions)

type ReplyOptions struct {
	Header Header
}

// WithReplyHeader 设置回复消息头。
func WithReplyHeader(h Header) ReplyOption {
	return func(o *ReplyOptions) {
		o.Header = h
	}
}
