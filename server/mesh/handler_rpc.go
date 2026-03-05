package mesh

import (
	"fmt"
	"reflect"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/encoding/proto"
)

type RpcMessageHandler func(*Mesh, *broker.Message) ([]byte, string, int)

// WrapRpc 路由处理函数包装器
// 统一处理request-reply消息,处理系统事件,自动解析业务参数 payload
// handler 返回:
//   - []byte: 业务数据
//   - string: 错误提示
//   - int: 业务状态码(200表示成功, 其它表示失败)
func WrapRpc[T any](handler func(*RpcContext, *T) ([]byte, string, int)) RpcMessageHandler {
	return func(m *Mesh, msg *broker.Message) ([]byte, string, int) {

		ctx := newRpcContext(m, msg)
		defer ctx.release()

		if len(msg.Data) == 0 {
			return handler(ctx, nil)
		}
		var payload T
		if err := proto.Unmarshal(msg.Data, &payload); err != nil {
			return nil, "unmarshal payload error", 500
		}
		return handler(ctx, &payload)
	}
}

// adaptRpcMessageHandler 将不同签名的 request handler 统一适配为 RpcMessageHandler
// 原理:
// 1) 若本身就是 RpcMessageHandler，直接返回
// 2) 若是 func(*RpcContext, *T) ([]byte, string, int)，使用反射校验签名后包装
// 3) 包装函数内统一完成 payload 反序列化、调用业务 handler、回包
func adaptRpcMessageHandler(handler any) (RpcMessageHandler, error) {
	if handler == nil {
		return nil, fmt.Errorf("mesh: request-reply handler is nil")
	}

	if mh, ok := handler.(RpcMessageHandler); ok {
		return mh, nil
	}

	rv := reflect.ValueOf(handler)
	rt := rv.Type()
	if rt.Kind() != reflect.Func {
		return nil, fmt.Errorf("mesh: unsupported route handler type %T", handler)
	}
	if rt.NumIn() != 2 || rt.NumOut() != 3 {
		return nil, fmt.Errorf("mesh: handler must be func(*RpcContext,*T)([]byte,string,int) or RpcMessageHandler, got %s", rt.String())
	}

	ctxType := reflect.TypeOf((*RpcContext)(nil))
	if rt.In(0) != ctxType {
		return nil, fmt.Errorf("mesh: handler first arg must be *mesh.RpcContext, got %s", rt.In(0).String())
	}
	argType := rt.In(1)
	if argType.Kind() != reflect.Ptr {
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
		ctx := newRpcContext(m, msg)
		defer ctx.release()

		callArg := reflect.Zero(argType)
		if len(msg.Data) > 0 {
			callArg = reflect.New(argType.Elem())
			if err := proto.Unmarshal(msg.Data, callArg.Interface()); err != nil {
				return nil, "unmarshal payload error", 500
			}
		}

		out := rv.Call([]reflect.Value{reflect.ValueOf(ctx), callArg})
		data := out[0].Interface().([]byte)
		tip := out[1].Interface().(string)
		code := int(out[2].Int())
		return data, tip, code
	}, nil
}
