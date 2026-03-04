package mesh

import (
	"fmt"
	"reflect"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/encoding/proto"
)

type RequestMessageHandler func(*Mesh, *broker.Message)

// WrapRequest 路由处理函数包装器
// 统一处理request-reply消息,处理系统事件,自动解析业务参数 payload
func WrapRequest[T any](handler func(*RequestContext, *T)) RequestMessageHandler {
	return func(m *Mesh, msg *broker.Message) {

		ctx := newRequestContext(m, msg)
		defer ctx.release()

		if len(msg.Data) == 0 {
			handler(ctx, nil)
			return
		}
		var payload T
		if err := proto.Unmarshal(msg.Data, &payload); err != nil {
			log.Errorf("mesh request-reply unmarshal payload error: %v", err)
			return
		}
		handler(ctx, &payload)
	}
}

// adaptRequestMessageHandler 将不同签名的 request handler 统一适配为 RequestMessageHandler
// 原理:
// 1) 若本身就是 RequestMessageHandler，直接返回
// 2) 若是 func(*RequestContext, *T)，使用反射校验签名后包装
// 3) 包装函数内统一完成 payload 反序列化，再调用业务 handler
func adaptRequestMessageHandler(handler any) (RequestMessageHandler, error) {
	if handler == nil {
		return nil, fmt.Errorf("mesh: request-reply handler is nil")
	}

	if mh, ok := handler.(RequestMessageHandler); ok {
		return mh, nil
	}

	rv := reflect.ValueOf(handler)
	rt := rv.Type()
	if rt.Kind() != reflect.Func {
		return nil, fmt.Errorf("mesh: unsupported route handler type %T", handler)
	}
	if rt.NumIn() != 2 || rt.NumOut() != 0 {
		return nil, fmt.Errorf("mesh: handler must be func(*RequestContext,*T) or RequestMessageHandler, got %s", rt.String())
	}

	ctxType := reflect.TypeOf((*RequestContext)(nil))
	if rt.In(0) != ctxType {
		return nil, fmt.Errorf("mesh: handler first arg must be *mesh.RequestContext, got %s", rt.In(0).String())
	}
	argType := rt.In(1)
	if argType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("mesh: handler second arg must be pointer type, got %s", argType.String())
	}
	return func(m *Mesh, msg *broker.Message) {
		ctx := newRequestContext(m, msg)
		defer ctx.release()

		callArg := reflect.Zero(argType)
		if len(msg.Data) > 0 {
			callArg = reflect.New(argType.Elem())
			if err := proto.Unmarshal(msg.Data, callArg.Interface()); err != nil {
				log.Errorf("mesh request-reply unmarshal payload error: %v", err)
				return
			}
		}

		rv.Call([]reflect.Value{reflect.ValueOf(ctx), callArg})
	}, nil
}
