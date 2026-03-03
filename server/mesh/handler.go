package mesh

import (
	"fmt"
	"reflect"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/encoding/proto"
	"github.com/byteweap/wukong/internal/envelope"
)

type MessageHandler func(*Mesh, *broker.Message) error

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// Wrap 路由处理函数包装器
// 统一处理网关消息,处理系统事件,自动解析业务参数 payload
func Wrap[I any](handler func(*Context, *I) error) MessageHandler {
	return func(m *Mesh, msg *broker.Message) error {

		envy := &envelope.Gate2MeshEnvelope{}
		if err := proto.Unmarshal(msg.Data, envy); err != nil {
			return err
		}
		switch envy.Event {
		case envelope.Event_ONLINE:
			if m.onlineHandler != nil {
				m.onlineHandler(envy.Uid)
			}
		case envelope.Event_OFFLINE:
			if m.offlineHandler != nil {
				m.offlineHandler(envy.Uid)
			}
		case envelope.Event_RECONNECT:
			if m.reconnectHandler != nil {
				m.reconnectHandler(envy.Uid)
			}
		case envelope.Event_Business:
			ctx := newContext(m, msg, envy)
			meta := envy.GetMeta()
			if meta == nil || len(meta.GetPayload()) == 0 {
				return handler(ctx, nil)
			}
			var payload I
			if err := proto.Unmarshal(meta.GetPayload(), &payload); err != nil {
				return err
			}
			return handler(ctx, &payload)
		}
		return nil
	}
}

// adaptMessageHandler 适配消息处理函数
func adaptMessageHandler(handler any) (MessageHandler, error) {
	if handler == nil {
		return nil, fmt.Errorf("mesh: handler is nil")
	}

	if mh, ok := handler.(MessageHandler); ok {
		return mh, nil
	}

	rv := reflect.ValueOf(handler)
	rt := rv.Type()
	if rt.Kind() != reflect.Func {
		return nil, fmt.Errorf("mesh: unsupported route handler type %T", handler)
	}
	if rt.NumIn() != 2 || rt.NumOut() != 1 {
		return nil, fmt.Errorf("mesh: handler must be func(*Context,*T) error or MessageHandler, got %s", rt.String())
	}

	ctxType := reflect.TypeOf((*Context)(nil))
	if rt.In(0) != ctxType {
		return nil, fmt.Errorf("mesh: handler first arg must be *mesh.Context, got %s", rt.In(0).String())
	}
	argType := rt.In(1)
	if argType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("mesh: handler second arg must be pointer type, got %s", argType.String())
	}
	if !rt.Out(0).Implements(errorType) {
		return nil, fmt.Errorf("mesh: handler return type must be error, got %s", rt.Out(0).String())
	}

	return func(m *Mesh, msg *broker.Message) error {
		envy := &envelope.Gate2MeshEnvelope{}
		if err := proto.Unmarshal(msg.Data, envy); err != nil {
			return err
		}

		switch envy.Event {
		case envelope.Event_ONLINE:
			if m.onlineHandler != nil {
				m.onlineHandler(envy.Uid)
			}
			return nil
		case envelope.Event_OFFLINE:
			if m.offlineHandler != nil {
				m.offlineHandler(envy.Uid)
			}
			return nil
		case envelope.Event_RECONNECT:
			if m.reconnectHandler != nil {
				m.reconnectHandler(envy.Uid)
			}
			return nil
		case envelope.Event_Business:
			ctx := newContext(m, msg, envy)
			meta := envy.GetMeta()

			callArg := reflect.Zero(argType)
			if meta != nil && len(meta.GetPayload()) > 0 {
				callArg = reflect.New(argType.Elem())
				if err := proto.Unmarshal(meta.GetPayload(), callArg.Interface()); err != nil {
					return err
				}
			}

			out := rv.Call([]reflect.Value{reflect.ValueOf(ctx), callArg})
			if out[0].IsNil() {
				return nil
			}
			return out[0].Interface().(error)
		default:
			return nil
		}
	}, nil
}
