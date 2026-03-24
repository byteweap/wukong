package mesh

import (
	"bytes"
	"strings"
	"testing"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/encoding/proto"
	"github.com/byteweap/wukong/envelope"
)

// TestAdaptRequestAutoWrapPayload 验证自动适配能正确反序列化请求参数
func TestAdaptRequestAutoWrapPayload(t *testing.T) {
	m := New()

	called := false
	var gotCtx *RpcContext
	var gotReq *envelope.IMessage

	wantData := []byte("ok-data")
	h := mustAdaptRequestHandler(t, func(ctx *RpcContext, req *envelope.IMessage) ([]byte, string, int) {
		called = true
		gotCtx = ctx
		gotReq = req
		return wantData, "ok", 200
	})

	raw, err := proto.Marshal(&envelope.IMessage{
		Header: &envelope.Header{
			Seq: 11,
			Cmd: 88,
		},
		Service: "rpc",
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	gotData, gotTip, gotCode := h(m, &broker.Message{
		Subject: "svc.game.request",
		Reply:   "svc.game.reply",
		Header: broker.Header{
			"cmd":     []string{"88"},
			"version": []string{"1"},
		},
		Data: raw,
	})

	if !called {
		t.Fatalf("request handler not called")
	}
	if gotCtx == nil {
		t.Fatalf("request context should not be nil")
	}
	if gotReq == nil {
		t.Fatalf("request payload should not be nil")
	}
	if gotReq.GetService() != "rpc" || gotReq.GetHeader().GetCmd() != 88 || gotReq.GetHeader().GetSeq() != 11 {
		t.Fatalf("unexpected payload: %+v", gotReq)
	}
	if !bytes.Equal(gotData, wantData) || gotTip != "ok" || gotCode != 200 {
		t.Fatalf("unexpected return: data=%v tip=%s code=%d", gotData, gotTip, gotCode)
	}
}

// TestAdaptRequestAutoWrapEmptyPayloadPassNil 验证空 payload 会传入 nil 参数
func TestAdaptRequestAutoWrapEmptyPayloadPassNil(t *testing.T) {
	m := New()

	var gotReq *envelope.IMessage
	h := mustAdaptRequestHandler(t, func(_ *RpcContext, req *envelope.IMessage) ([]byte, string, int) {
		gotReq = req
		return nil, "ok", 200
	})

	_, _, _ = h(m, &broker.Message{Data: nil})
	if gotReq != nil {
		t.Fatalf("expected nil request for empty payload")
	}
}

// TestAdaptRequestCompatibleWithRequestMessageHandler 验证兼容直接传 RpcMessageHandler
func TestAdaptRequestCompatibleWithRequestMessageHandler(t *testing.T) {
	m := New()

	called := false
	h, err := adaptRpcMessageHandler(RpcMessageHandler(func(_ *Mesh, _ *broker.Message) ([]byte, string, int) {
		called = true
		return nil, "ok", 200
	}))
	if err != nil {
		t.Fatalf("adapt request handler: %v", err)
	}

	_, _, _ = h(m, &broker.Message{})
	if !called {
		t.Fatalf("request message handler not called")
	}
}

// TestAdaptRequestInvalidHandlerError 验证非法签名会返回错误
func TestAdaptRequestInvalidHandlerError(t *testing.T) {
	_, err := adaptRpcMessageHandler(func() error { return nil })
	if err == nil {
		t.Fatalf("expected error for invalid request handler")
	}
	if !strings.Contains(err.Error(), "func(*RpcContext,*T)([]byte,string,int)") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func mustAdaptRequestHandler(t *testing.T, handler any) RpcMessageHandler {
	t.Helper()
	h, err := adaptRpcMessageHandler(handler)
	if err != nil {
		t.Fatalf("adapt request handler: %v", err)
	}
	return h
}

// TestRequestRouteDispatchByHeader 验证 request-reply 根据 header 中的 cmd/version 分发
func TestRequestRouteDispatchByHeader(t *testing.T) {
	mb := &mockBroker{}
	m := New(Broker(mb))

	called := false
	var gotCtx *RpcContext
	var gotReq *envelope.IMessage

	wantData := []byte("route-ok")
	m.RpcRouteX("2001", "1", func(ctx *RpcContext, req *envelope.IMessage) ([]byte, string, int) {
		called = true
		gotCtx = ctx
		gotReq = req
		return wantData, "ok", 200
	})

	raw, err := proto.Marshal(&envelope.IMessage{
		Header: &envelope.Header{
			Seq: 21,
			Cmd: 2001,
		},
		Service: "mesh",
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	m.handlerRequestReplyMessage(&broker.Message{
		Subject: "svc.request",
		Reply:   "svc.reply",
		Header: broker.Header{
			"cmd":     []string{"2001"},
			"version": []string{"1"},
		},
		Data: raw,
	})

	if !called {
		t.Fatalf("request route handler not called")
	}
	if gotCtx == nil {
		t.Fatalf("request context should not be nil")
	}
	if gotReq == nil {
		t.Fatalf("request payload should not be nil")
	}

	if gotReq.GetService() != "mesh" || gotReq.GetHeader().GetCmd() != 2001 || gotReq.GetHeader().GetSeq() != 21 {
		t.Fatalf("unexpected payload: %+v", gotReq)
	}
	if mb.replyCalls != 1 {
		t.Fatalf("expected 1 reply call, got %d", mb.replyCalls)
	}
	if !bytes.Equal(mb.replyData, wantData) {
		t.Fatalf("unexpected reply data: %v", mb.replyData)
	}
	if mb.replyHdr.Get("code") != "200" || mb.replyHdr.Get("tip") != "ok" {
		t.Fatalf("unexpected reply header: %+v", mb.replyHdr)
	}
}

// TestRequestRouteDispatchEmptyPayloadPassNil 验证空 payload 分发时传入 nil 参数
func TestRequestRouteDispatchEmptyPayloadPassNil(t *testing.T) {
	mb := &mockBroker{}
	m := New(Broker(mb))

	var gotReq *envelope.IMessage
	m.RpcRouteX("2002", "1", func(_ *RpcContext, req *envelope.IMessage) ([]byte, string, int) {
		gotReq = req
		return nil, "ok", 200
	})

	m.handlerRequestReplyMessage(&broker.Message{
		Reply: "svc.reply",
		Header: broker.Header{
			"cmd":     []string{"2002"},
			"version": []string{"1"},
		},
		Data: nil,
	})

	if gotReq != nil {
		t.Fatalf("expected nil request for empty payload")
	}
}

func TestRequestRouteErrorReply(t *testing.T) {
	mb := &mockBroker{}
	m := New(Broker(mb))

	m.RpcRouteX("2003", "1", func(_ *RpcContext, _ *envelope.IMessage) ([]byte, string, int) {
		return []byte("ignored"), "bad request", 400
	})

	m.handlerRequestReplyMessage(&broker.Message{
		Reply: "svc.reply",
		Header: broker.Header{
			"cmd":     []string{"2003"},
			"version": []string{"1"},
		},
		Data: nil,
	})

	if mb.replyCalls != 1 {
		t.Fatalf("expected 1 reply call, got %d", mb.replyCalls)
	}
	if mb.replyData != nil {
		t.Fatalf("error reply data should be nil")
	}
	if mb.replyHdr.Get("code") != "400" || mb.replyHdr.Get("tip") != "bad request" {
		t.Fatalf("unexpected reply header: %+v", mb.replyHdr)
	}
}

// TestRequestRouteInvalidHandlerPanic 验证 RpcRouteX 对非法签名会 panic
func TestRequestRouteInvalidHandlerPanic(t *testing.T) {
	m := New()

	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic for invalid request handler")
		}
	}()

	m.RpcRouteX("9999", "1", func() error { return nil })
}
