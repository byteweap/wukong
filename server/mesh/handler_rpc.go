package mesh

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/byteweap/meta/component/broker"
	"github.com/byteweap/meta/encoding/proto"
)

type RPCMessageHandler func(*Mesh, *broker.Message) ([]byte, string, int)

const unmarshalPayloadError = "unmarshal payload error"

// WrapRPC 路由处理函数包装器
// 统一处理request-reply消息,处理系统事件,自动解析业务参数 payload
// handler 返回:
//   - []byte: 业务数据
//   - string: 错误提示
//   - int: 业务状态码(200表示成功, 其它表示失败)
func WrapRPC[T any](handler func(*RPCContext, *T) ([]byte, string, int)) RPCMessageHandler {
	return func(m *Mesh, msg *broker.Message) ([]byte, string, int) {
		ctx := newRPCContext(m, msg)
		defer ctx.release()

		if len(msg.Data) == 0 {
			return handler(ctx, nil)
		}
		var payload T
		if err := proto.Unmarshal(msg.Data, &payload); err != nil {
			return nil, unmarshalPayloadError, http.StatusInternalServerError
		}
		return handler(ctx, &payload)
	}
}

// adaptRPCMessageHandler 将不同签名的 request handler 统一适配为 RPCMessageHandler
// 原理:
// 1) 若本身就是 RPCMessageHandler，直接返回
// 2) 若是 func(*RPCContext, *T) ([]byte, string, int)，使用反射校验签名后包装
// 3) 包装函数内统一完成 payload 反序列化、调用业务 handler、回包
func adaptRPCMessageHandler(handler any) (RPCMessageHandler, error) {
	if handler == nil {
		return nil, fmt.Errorf("mesh: request-reply handler is nil")
	}

	if mh, ok := handler.(RPCMessageHandler); ok {
		return mh, nil
	}

	rv := reflect.ValueOf(handler)
	rt := rv.Type()
	if rt.Kind() != reflect.Func {
		return nil, fmt.Errorf("mesh: unsupported route handler type %T", handler)
	}
	if rt.NumIn() != 2 || rt.NumOut() != 3 {
		return nil, fmt.Errorf("mesh: handler must be func(*RPCContext,*T)([]byte,string,int) or RPCMessageHandler, got %s", rt.String())
	}

	ctxType := reflect.TypeOf((*RPCContext)(nil))
	if rt.In(0) != ctxType {
		return nil, fmt.Errorf("mesh: handler first arg must be *mesh.RPCContext, got %s", rt.In(0).String())
	}
	argType := rt.In(1)
	if argType.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("mesh: handler second arg must be pointer type, got %s", argType.String())
	}
	if rt.Out(0) != reflect.TypeOf([]byte(nil)) {
		return nil, fmt.Errorf("mesh: handler first return must be []byte, got %s", rt.Out(0).String())
	}
	if rt.Out(1).Kind() != reflect.String {
		return nil, fmt.Errorf("mesh: handler second return must be string, got %s", rt.Out(1).String())
	}
	if rt.Out(2).Kind() != reflect.Int {
		return nil, fmt.Errorf("mesh: handler third return must be int, got %s", rt.Out(2).String())
	}

	return func(m *Mesh, msg *broker.Message) ([]byte, string, int) {
		ctx := newRPCContext(m, msg)
		defer ctx.release()

		callArg := reflect.Zero(argType)
		if len(msg.Data) > 0 {
			callArg = reflect.New(argType.Elem())
			if err := proto.Unmarshal(msg.Data, callArg.Interface()); err != nil {
				return nil, unmarshalPayloadError, http.StatusInternalServerError
			}
		}

		out := rv.Call([]reflect.Value{reflect.ValueOf(ctx), callArg})
		data := out[0].Interface().([]byte)
		tip := out[1].Interface().(string)
		code := int(out[2].Int())
		return data, tip, code
	}, nil
}
