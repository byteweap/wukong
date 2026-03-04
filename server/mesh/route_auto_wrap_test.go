package mesh

import (
	"testing"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/encoding/proto"
	"github.com/byteweap/wukong/internal/envelope"
)

// TestRouteAutoWrapBusinessPayload 验证自动 Wrap 能正确反序列化业务参数
func TestRouteAutoWrapBusinessPayload(t *testing.T) {
	m := New()

	called := false
	var gotCtx *Context
	var gotReq *envelope.Envelope

	m.Route(1001, 1, func(ctx *Context, req *envelope.Envelope) {
		called = true
		gotCtx = ctx
		gotReq = req
	})

	raw := mustBusinessMessage(t, &envelope.Envelope{
		Seq: 99,
		App: "game",
		Cmd: 1001,
	})

	h := mustLoadRouteHandler(t, m, 1001, 1)
	invokeRouteHandler(t, h, m, &broker.Message{Data: raw}, raw)

	if !called {
		t.Fatalf("business handler not called")
	}
	if gotCtx == nil {
		t.Fatalf("context should not be nil")
	}
	if gotReq == nil {
		t.Fatalf("request should not be nil")
	}
	if gotReq.GetApp() != "game" || gotReq.GetCmd() != 1001 || gotReq.GetSeq() != 99 {
		t.Fatalf("unexpected payload: %+v", gotReq)
	}
}

// TestRouteAutoWrapEmptyPayloadPassNil 验证空 payload 会传入 nil 参数
func TestRouteAutoWrapEmptyPayloadPassNil(t *testing.T) {
	m := New()

	var gotReq *envelope.Envelope
	m.Route(1002, 1, func(_ *Context, req *envelope.Envelope) {
		gotReq = req
	})

	raw, err := proto.Marshal(&envelope.Gate2MeshEnvelope{
		Event: envelope.Event_Business,
		Meta: &envelope.Envelope{
			Cmd: 1002,
		},
	})
	if err != nil {
		t.Fatalf("marshal gate envelope: %v", err)
	}

	h := mustLoadRouteHandler(t, m, 1002, 1)
	invokeRouteHandler(t, h, m, &broker.Message{Data: raw}, raw)
	if gotReq != nil {
		t.Fatalf("expected nil request for empty payload")
	}
}

// TestRouteCompatibleWithMessageHandler 验证兼容直接注册 MessageHandler
func TestRouteCompatibleWithMessageHandler(t *testing.T) {
	m := New()

	called := false
	m.Route(1003, 1, MessageHandler(func(_ *Mesh, _ *broker.Message, _ *envelope.Gate2MeshEnvelope) {
		called = true
	}))

	h := mustLoadRouteHandler(t, m, 1003, 1)
	h(m, &broker.Message{}, &envelope.Gate2MeshEnvelope{})
	if !called {
		t.Fatalf("message handler not called")
	}
}

// TestRouteOnlineEventWithoutCallback 验证系统事件不会误触发业务处理器
func TestRouteOnlineEventWithoutCallback(t *testing.T) {
	m := New()

	called := false
	m.Route(1004, 1, func(_ *Context, _ *envelope.Envelope) {
		called = true
	})

	raw, err := proto.Marshal(&envelope.Gate2MeshEnvelope{
		Event: envelope.Event_ONLINE,
		Uid:   7,
	})
	if err != nil {
		t.Fatalf("marshal gate envelope: %v", err)
	}

	h := mustLoadRouteHandler(t, m, 1004, 1)
	invokeRouteHandler(t, h, m, &broker.Message{Data: raw}, raw)
	if called {
		t.Fatalf("business handler should not be called for online event")
	}
}

// TestRouteInvalidHandlerPanic 验证非法签名会触发 panic
func TestRouteInvalidHandlerPanic(t *testing.T) {
	m := New()

	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic for invalid handler")
		}
	}()

	m.Route(9999, 1, func() error { return nil })
}

// mustBusinessMessage 构造业务消息封包
func mustBusinessMessage(t *testing.T, payload *envelope.Envelope) []byte {
	t.Helper()

	p, err := proto.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	raw, err := proto.Marshal(&envelope.Gate2MeshEnvelope{
		Event: envelope.Event_Business,
		Meta: &envelope.Envelope{
			Cmd:     payload.GetCmd(),
			Payload: p,
		},
	})
	if err != nil {
		t.Fatalf("marshal gate envelope: %v", err)
	}
	return raw
}

// mustLoadRouteHandler 读取已注册路由处理器
func mustLoadRouteHandler(t *testing.T, m *Mesh, cmd, version uint32) MessageHandler {
	t.Helper()

	v, ok := m.routes.Load(routeKey(cmd, version))
	if !ok {
		t.Fatalf("route handler not found")
	}
	h, ok := v.(MessageHandler)
	if !ok {
		t.Fatalf("stored handler has unexpected type %T", v)
	}
	return h
}

func invokeRouteHandler(t *testing.T, h MessageHandler, m *Mesh, msg *broker.Message, raw []byte) {
	t.Helper()

	envy := &envelope.Gate2MeshEnvelope{}
	if err := proto.Unmarshal(raw, envy); err != nil {
		t.Fatalf("unmarshal gate envelope: %v", err)
	}
	h(m, msg, envy)
}
