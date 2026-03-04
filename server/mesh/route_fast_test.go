package mesh

import (
	"testing"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/encoding/proto"
	"github.com/byteweap/wukong/internal/envelope"
)

func TestRouteFastWithWrap(t *testing.T) {
	m := New()

	called := false
	m.Route(3001, 1, Wrap(func(_ *Context, req *envelope.Envelope) {
		if req != nil && req.GetCmd() == 3001 {
			called = true
		}
	}))

	raw := mustBusinessMessage(t, &envelope.Envelope{
		Cmd: 3001,
	})

	h := mustLoadRouteHandler(t, m, 3001, 1)
	invokeRouteHandler(t, h, m, &broker.Message{Data: raw}, raw)
	if !called {
		t.Fatalf("route fast handler not called")
	}
}

func TestRequestRouteFastWithWrapRequest(t *testing.T) {
	m := New()

	called := false
	m.RequestRoute("findUser", "v1", WrapRequest(func(_ *RequestContext, req *envelope.Envelope) {
		if req != nil && req.GetApp() == "mesh" {
			called = true
		}
	}))

	raw, err := proto.Marshal(&envelope.Envelope{App: "mesh"})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	m.handlerRequestReplyMessage(&broker.Message{
		Reply: "svc.reply",
		Header: broker.Header{
			"cmd":     []string{"findUser"},
			"version": []string{"v1"},
		},
		Data: raw,
	})

	if !called {
		t.Fatalf("request route fast handler not called")
	}
}
