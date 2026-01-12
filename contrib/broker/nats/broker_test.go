package nats

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	natssrv "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"

	"github.com/byteweap/wukong/component/broker"
)

func runNatsServer(t *testing.T) *natssrv.Server {
	t.Helper()

	opts := &natssrv.Options{
		Host:   "127.0.0.1",
		Port:   -1, // random available port
		NoLog:  true,
		NoSigs: true,
	}
	s, err := natssrv.NewServer(opts)
	require.NoError(t, err)

	go s.Start()
	require.True(t, s.ReadyForConnections(5*time.Second), "nats server not ready")

	t.Cleanup(func() {
		s.Shutdown()
	})
	return s
}

func TestPublishSubscribe_HeaderRoundTrip(t *testing.T) {
	s := runNatsServer(t)

	b, err := New(WithURLs(s.ClientURL()))
	require.NoError(t, err)
	t.Cleanup(b.Close)

	var (
		subject = "t.pubsub.v1"
		data    = []byte("hello")
	)

	got := make(chan *broker.Message, 1)
	_, err = b.Sub(context.Background(), subject, func(_ context.Context, msg *broker.Message) {
		got <- msg
	})
	require.NoError(t, err)

	h := broker.Header{
		"X-Trace-Id": {"abc"},
	}
	err = b.Pub(context.Background(), subject, data, broker.WithHeader(h))
	require.NoError(t, err)

	select {
	case msg := <-got:

		require.Equal(t, subject, msg.Subject)
		require.Equal(t, data, msg.Data)
		require.Equal(t, h, msg.Header)

		t.Logf("msg.Header: %+v", msg.Header)
		t.Logf("msg.Data: %s", string(msg.Data))

	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for message")
	}
}

func TestQueueSubscribe_ExactlyOnce(t *testing.T) {
	s := runNatsServer(t)

	b, err := New(WithURLs(s.ClientURL()))
	require.NoError(t, err)
	t.Cleanup(b.Close)

	const (
		subject = "t.queue.v1"
		queue   = "q1"
		N       = 200
	)

	var c1, c2 int32
	ctx := context.Background()

	_, err = b.Sub(ctx, subject, func(_ context.Context, _ *broker.Message) {
		atomic.AddInt32(&c1, 1)
	}, broker.WithQueue(queue))
	require.NoError(t, err)

	_, err = b.Sub(ctx, subject, func(_ context.Context, _ *broker.Message) {
		atomic.AddInt32(&c2, 1)
	}, broker.WithQueue(queue))
	require.NoError(t, err)

	// Give subscriptions a tiny moment to propagate.
	time.Sleep(50 * time.Millisecond)

	for i := 0; i < N; i++ {
		require.NoError(t, b.Pub(ctx, subject, []byte("x")))
	}

	require.Eventually(t, func() bool {
		return int(atomic.LoadInt32(&c1)+atomic.LoadInt32(&c2)) == N
	}, 2*time.Second, 10*time.Millisecond, "expected exactly %d messages consumed once across queue group", N)
}

func TestRequestReply(t *testing.T) {
	s := runNatsServer(t)

	b, err := New(WithURLs(s.ClientURL()))
	require.NoError(t, err)
	t.Cleanup(b.Close)

	var (
		subject = "t.req.v1"
		ping    = []byte("ping")
		pong    = []byte("pong")
	)

	_, err = b.Sub(context.Background(), subject, func(ctx context.Context, msg *broker.Message) {
		// 使用 Reply 方法回复请求（更语义化）
		_ = b.Reply(ctx, msg, pong, broker.WithReplyHeader(broker.Header{
			"X-From": {"server"},
		}))
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	resp, err := b.Request(
		ctx,
		subject,
		ping,
		broker.WithRequestHeader(broker.Header{
			"X-From": {"client"},
		}))
	require.NoError(t, err)
	require.Equal(t, pong, resp.Data)
	require.Equal(t, broker.Header{"X-From": {"server"}}, resp.Header)
}

func TestSubscribe_ContextCancelAutoUnsubscribe(t *testing.T) {
	s := runNatsServer(t)

	b, err := New(WithURLs(s.ClientURL()))
	require.NoError(t, err)
	t.Cleanup(b.Close)

	var (
		subject = "t.cancel.v1"
	)

	ctx, cancel := context.WithCancel(context.Background())
	called := make(chan struct{}, 1)

	sub, err := b.Sub(ctx, subject, func(_ context.Context, _ *broker.Message) {
		called <- struct{}{}
	})
	require.NoError(t, err)

	// cancel -> broker should unsubscribe
	cancel()

	// assert underlying nats subscription becomes invalid eventually
	ns := sub.(*subscription).sub
	require.Eventually(t, func() bool {
		return !ns.IsValid()
	}, 2*time.Second, 10*time.Millisecond)

	// publishing afterwards should not call handler
	_ = b.Pub(context.Background(), subject, []byte("x"))
	select {
	case <-called:
		t.Fatal("handler should not be called after ctx cancel unsubscribe")
	case <-time.After(150 * time.Millisecond):
		// ok
	}
}

func TestContextErrFastFail(t *testing.T) {
	s := runNatsServer(t)

	b, err := New(WithURLs(s.ClientURL()))
	require.NoError(t, err)
	t.Cleanup(b.Close)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	require.ErrorIs(t, b.Pub(ctx, "t.any", []byte("x")), context.Canceled)

	_, err = b.Request(ctx, "t.any", []byte("x"))
	require.ErrorIs(t, err, context.Canceled)
}

func TestReply(t *testing.T) {
	s := runNatsServer(t)

	b, err := New(WithURLs(s.ClientURL()))
	require.NoError(t, err)
	t.Cleanup(b.Close)

	// 测试：正常 reply
	var gotReply bool
	_, err = b.Sub(context.Background(), "t.reply.v1", func(ctx context.Context, msg *broker.Message) {
		err := b.Reply(ctx, msg, []byte("ok"), broker.WithReplyHeader(broker.Header{
			"X-Status": {"success"},
		}))
		require.NoError(t, err)
		gotReply = true
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	resp, err := b.Request(ctx, "t.reply.v1", []byte("test"))
	require.NoError(t, err)
	require.Equal(t, []byte("ok"), resp.Data)
	require.Equal(t, broker.Header{"X-Status": {"success"}}, resp.Header)
	require.Eventually(t, func() bool { return gotReply }, 100*time.Millisecond, 10*time.Millisecond)

	// 测试：msg.Reply 为空应该返回错误
	msgWithoutReply := &broker.Message{Subject: "test", Reply: "", Data: []byte("x")}
	err = b.Reply(context.Background(), msgWithoutReply, []byte("ok"))
	require.Error(t, err)
	require.True(t, errors.Is(err, nats.ErrInvalidMsg))

	// 测试：msg 为 nil 应该返回错误
	err = b.Reply(context.Background(), nil, []byte("ok"))
	require.Error(t, err)
	require.True(t, errors.Is(err, nats.ErrInvalidMsg))

	// 测试：ctx 取消应该快速失败
	ctx2, cancel2 := context.WithCancel(context.Background())
	cancel2()
	msg := &broker.Message{Subject: "test", Reply: "reply.inbox", Data: []byte("x")}
	require.ErrorIs(t, b.Reply(ctx2, msg, []byte("ok")), context.Canceled)
}
