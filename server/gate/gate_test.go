package gate

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/olahol/melody"
	"github.com/stretchr/testify/require"

	"github.com/byteweap/meta"
	"github.com/byteweap/meta/component/broker"
	"github.com/byteweap/meta/component/locator"
	"github.com/byteweap/meta/component/registry"
	"github.com/byteweap/meta/component/selector"
)

type testAppInfo struct {
	id   string
	name string
}

func (a testAppInfo) ID() string                  { return a.id }
func (a testAppInfo) Name() string                { return a.name }
func (a testAppInfo) Version() string             { return "v1.0.0" }
func (a testAppInfo) Metadata() map[string]string { return map[string]string{} }
func (a testAppInfo) Endpoint() []string          { return nil }

type testSubscription struct{}

func (s *testSubscription) Unsub() error { return nil }
func (s *testSubscription) Close() error { return nil }

type testBroker struct {
	mu          sync.Mutex
	replyCalls  int
	replyData   []byte
	replyHeader broker.Header
}

func (b *testBroker) ID() string { return "test-broker" }

func (b *testBroker) Pub(context.Context, string, []byte, ...broker.PublishOption) error {
	return nil
}

func (b *testBroker) Sub(context.Context, string, broker.Handler, ...broker.SubscribeOption) (broker.Subscription, error) {
	return &testSubscription{}, nil
}

func (b *testBroker) Request(context.Context, string, []byte, ...broker.RequestOption) (*broker.Message, error) {
	return &broker.Message{}, nil
}

func (b *testBroker) Reply(_ context.Context, _ *broker.Message, data []byte, opts ...broker.ReplyOption) error {
	replyOpts := &broker.ReplyOptions{}
	for _, opt := range opts {
		opt(replyOpts)
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.replyCalls++
	b.replyData = append([]byte(nil), data...)
	b.replyHeader = replyOpts.Header
	return nil
}

func (b *testBroker) Close() error { return nil }

type testLocator struct {
	mu          sync.Mutex
	bindErr     error
	bindCalls   int
	unbindCalls int
}

var _ locator.Locator = (*testLocator)(nil)

func (l *testLocator) ID() string { return "test-locator" }

func (l *testLocator) AllNodes(context.Context, int64) (map[string]string, error) {
	return map[string]string{}, nil
}

func (l *testLocator) Node(context.Context, int64, string) (string, error) {
	return "", nil
}

func (l *testLocator) Bind(context.Context, int64, string, string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.bindCalls++
	return l.bindErr
}

func (l *testLocator) UnBind(context.Context, int64, string, string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.unbindCalls++
	return nil
}

func (l *testLocator) Close() error { return nil }

type testWatcher struct {
	stopCalls int
	nextErr   error
}

func (w *testWatcher) Next() ([]*registry.ServiceInstance, error) {
	if w.nextErr != nil {
		return nil, w.nextErr
	}
	return []*registry.ServiceInstance{}, nil
}

func (w *testWatcher) Stop() error {
	w.stopCalls++
	return nil
}

type testRegistry struct {
	mu         sync.Mutex
	watchCalls int
	watchErrs  []error
	watcher    registry.Watcher
}

var _ registry.Registry = (*testRegistry)(nil)

func (r *testRegistry) ID() string { return "test-registry" }

func (r *testRegistry) Register(context.Context, *registry.ServiceInstance) error { return nil }

func (r *testRegistry) Deregister(context.Context, *registry.ServiceInstance) error { return nil }

func (r *testRegistry) GetService(context.Context, string) ([]*registry.ServiceInstance, error) {
	return nil, nil
}

func (r *testRegistry) Watch(context.Context, string) (registry.Watcher, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.watchCalls++
	idx := r.watchCalls - 1
	if idx < len(r.watchErrs) && r.watchErrs[idx] != nil {
		return nil, r.watchErrs[idx]
	}
	return r.watcher, nil
}

type testSelector struct {
	mu    sync.Mutex
	nodes []selector.Node
}

func (s *testSelector) Select(string, ...selector.Filter) (selector.Node, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.nodes) == 0 {
		return nil, selector.ErrNoAvailableNode
	}
	return s.nodes[0], nil
}

func (s *testSelector) Update(nodes []selector.Node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodes = append([]selector.Node(nil), nodes...)
}

func (s *testSelector) Nodes() []selector.Node {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]selector.Node(nil), s.nodes...)
}

func TestGateValidate(t *testing.T) {
	g := New()
	require.Error(t, g.validate())

	g = New(
		Locator(&testLocator{}),
		Broker(&testBroker{}),
		Discovery(&testRegistry{}),
		SelectorFunc(func() selector.Selector { return &testSelector{} }),
	)
	require.NoError(t, g.validate())
}

func TestEnsureRetriesAfterWatchFailure(t *testing.T) {
	watcher := &testWatcher{nextErr: errors.New("stop")}
	discovery := &testRegistry{
		watchErrs: []error{errors.New("watch failed"), nil},
		watcher:   watcher,
	}

	g := New(
		Discovery(discovery),
		SelectorFunc(func() selector.Selector { return &testSelector{} }),
	)
	g.ctx = context.Background()

	_, err := g.ensure("match")
	require.EqualError(t, err, "watch failed")
	require.Nil(t, g.selectors["match"])
	require.Nil(t, g.watchers["match"])

	sel, err := g.ensure("match")
	require.NoError(t, err)
	require.NotNil(t, sel)
	require.NotNil(t, g.selectors["match"])
	require.NotNil(t, g.watchers["match"])
	require.Equal(t, 2, discovery.watchCalls)
}

func TestHandleConnectRollsBackSessionWhenBindFails(t *testing.T) {
	loc := &testLocator{bindErr: errors.New("bind failed")}
	g := New(
		Addr("127.0.0.1:0"),
		Path("/ws"),
		Locator(loc),
		Broker(&testBroker{}),
		Discovery(&testRegistry{}),
		SelectorFunc(func() selector.Selector { return &testSelector{} }),
	)

	ctx := meta.NewContext(context.Background(), testAppInfo{id: "gate-1", name: "gate"})
	require.NoError(t, g.setup("gate", "gate-1", ctx))
	defer func() {
		if g.ws != nil {
			_ = g.ws.Close()
		}
		if g.ln != nil {
			_ = g.ln.Close()
		}
	}()

	ts := httptest.NewServer(g.Handler)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws?uid=42"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	require.Eventually(t, func() bool {
		_, ok := g.sessions.get(42)
		return !ok
	}, time.Second, 20*time.Millisecond)

	require.Eventually(t, func() bool {
		_ = conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
		_, _, readErr := conn.ReadMessage()
		return readErr != nil
	}, 2*time.Second, 20*time.Millisecond)

	loc.mu.Lock()
	defer loc.mu.Unlock()
	require.Equal(t, 1, loc.bindCalls)
	require.Equal(t, 0, loc.unbindCalls)
}

func TestHandleRequestReplyMessageReturnsNotImplemented(t *testing.T) {
	bro := &testBroker{}
	g := New(Broker(bro))
	g.ctx = context.Background()

	g.handleRequestReplyMessage(&broker.Message{
		Reply: "reply.to.gate",
		Header: broker.Header{
			"trace": []string{"123"},
		},
	})

	bro.mu.Lock()
	defer bro.mu.Unlock()
	require.Equal(t, 1, bro.replyCalls)
	require.Nil(t, bro.replyData)
	require.Equal(t, "501", bro.replyHeader.Get("code"))
	require.Equal(t, "gate request-reply is not implemented", bro.replyHeader.Get("tip"))
	require.Equal(t, "123", bro.replyHeader.Get("trace"))
}

func TestHandleDisconnectIgnoresStaleSession(t *testing.T) {
	loc := &testLocator{}
	g := New(Locator(loc))
	g.ctx = context.Background()
	g.appName = "gate"
	g.appID = "gate-1"

	current := &melody.Session{Keys: map[string]any{"uid": int64(7)}}
	stale := &melody.Session{Keys: map[string]any{"uid": int64(7)}}
	g.sessions.register(7, current)

	g.handleDisconnect(stale)

	session, ok := g.sessions.get(7)
	require.True(t, ok)
	require.Same(t, current, session)

	loc.mu.Lock()
	defer loc.mu.Unlock()
	require.Equal(t, 0, loc.unbindCalls)
}
