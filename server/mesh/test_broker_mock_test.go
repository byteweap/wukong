package mesh

import (
	"context"

	"github.com/byteweap/wukong/component/broker"
)

type mockSubscription struct{}

func (s *mockSubscription) Unsub() error { return nil }
func (s *mockSubscription) Close() error { return nil }

type mockBroker struct {
	replyCalls int
	replyData  []byte
	replyHdr   broker.Header
}

func (b *mockBroker) ID() string { return "mock" }

func (b *mockBroker) Pub(ctx context.Context, subject string, data []byte, opts ...broker.PublishOption) error {
	return nil
}

func (b *mockBroker) Sub(ctx context.Context, subject string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscription, error) {
	return &mockSubscription{}, nil
}

func (b *mockBroker) Request(ctx context.Context, subject string, data []byte, opts ...broker.RequestOption) (*broker.Message, error) {
	return &broker.Message{}, nil
}

func (b *mockBroker) Reply(ctx context.Context, msg *broker.Message, data []byte, opts ...broker.ReplyOption) error {
	replyOpt := &broker.ReplyOptions{}
	for _, opt := range opts {
		opt(replyOpt)
	}
	b.replyCalls++
	b.replyData = data
	b.replyHdr = replyOpt.Header
	return nil
}

func (b *mockBroker) Close() error { return nil }
