package mesh

import (
	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/encoding/proto"
)

type MessageHandler func(*Mesh, *broker.Message) error

func Wrap[I any](handler func(*Context, *I) error) MessageHandler {
	return func(m *Mesh, msg *broker.Message) error {
		if len(msg.Data) > 0 {
			var input I
			if err := proto.Unmarshal(msg.Data, &input); err != nil {
				return err
			}
			return handler(newContext(m, msg), &input)
		}
		return handler(newContext(m, msg), nil)
	}
}
