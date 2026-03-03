package mesh

import (
	"testing"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/encoding/proto"
	"github.com/byteweap/wukong/internal/envelope"
)

func TestRouteAutoWrapBusinessPayload(t *testing.T) {
	m := New()

	called := false
	var gotCtx *Context
	var gotReq *envelope.Envelope

	m.Route(1001, 1, func(ctx *Context, req *envelope.Envelope) error {
		called = true
		gotCtx = ctx
		gotReq = req
		return nil
	})

	raw := mustBusinessMessage(t, &envelope.Envelope{
		Seq: 99,
		App: "game",
		Cmd: 1001,
	})

	h := mustLoadRouteHandler(t, m, 1001, 1)
	if err := h(m, &broker.Message{Data: raw}); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

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

func TestRouteAutoWrapEmptyPayloadPassNil(t *testing.T) {
	m := New()

	var gotReq *envelope.Envelope
	m.Route(1002, 1, func(_ *Context, req *envelope.Envelope) error {
		gotReq = req
		return nil
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
	if err := h(m, &broker.Message{Data: raw}); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if gotReq != nil {
		t.Fatalf("expected nil request for empty payload")
	}
}

func TestRouteCompatibleWithMessageHandler(t *testing.T) {
	m := New()

	called := false
	m.Route(1003, 1, MessageHandler(func(_ *Mesh, _ *broker.Message) error {
		called = true
		return nil
	}))

	h := mustLoadRouteHandler(t, m, 1003, 1)
	if err := h(m, &broker.Message{}); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if !called {
		t.Fatalf("message handler not called")
	}
}

func TestRouteOnlineEventWithoutCallback(t *testing.T) {
	m := New()

	called := false
	m.Route(1004, 1, func(_ *Context, _ *envelope.Envelope) error {
		called = true
		return nil
	})

	raw, err := proto.Marshal(&envelope.Gate2MeshEnvelope{
		Event: envelope.Event_ONLINE,
		Uid:   7,
	})
	if err != nil {
		t.Fatalf("marshal gate envelope: %v", err)
	}

	h := mustLoadRouteHandler(t, m, 1004, 1)
	if err := h(m, &broker.Message{Data: raw}); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if called {
		t.Fatalf("business handler should not be called for online event")
	}
}

func TestRouteInvalidHandlerPanic(t *testing.T) {
	m := New()

	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic for invalid handler")
		}
	}()

	m.Route(9999, 1, func() error { return nil })
}

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
