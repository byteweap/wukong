package broker

// PublishOption 发布选项 (预留).
type PublishOption func(*PublishOptions)

type PublishOptions struct {
	Header Header
	Reply  string // 可选的回复地址 (reply subject)，用于异步 request-reply
}

// SubscribeOption 订阅选项.
type SubscribeOption func(*SubscribeOptions)

type SubscribeOptions struct {
	Queue string // Queue Group（为空表示普通订阅/广播）
}

// RequestOption 请求选项.
type RequestOption func(*RequestOptions)

type RequestOptions struct {
	Header Header
}

// PubHeader 设置消息头 (Publish/Request).
func PubHeader(h Header) PublishOption {
	return func(o *PublishOptions) {
		o.Header = h
	}
}

// RequestHeader 设置消息头 (Request).
func RequestHeader(h Header) RequestOption {
	return func(o *RequestOptions) {
		o.Header = h
	}
}

// SubQueue 设置订阅的 Queue Group (用于水平扩展的竞争消费者).
func SubQueue(queue string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = queue
	}
}

// ReplyOption 回复选项.
type ReplyOption func(*ReplyOptions)

type ReplyOptions struct {
	Header Header
}

// ReplyHeader 设置回复消息头.
func ReplyHeader(h Header) ReplyOption {
	return func(o *ReplyOptions) {
		o.Header = h
	}
}

// PublishReply 设置发布消息的回复地址 (用于异步 request-reply).
func PublishReply(reply string) PublishOption {
	return func(o *PublishOptions) {
		o.Reply = reply
	}
}
