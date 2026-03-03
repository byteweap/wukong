package mesh

import (
	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/encoding/proto"
	"github.com/byteweap/wukong/internal/envelope"
)

type MessageHandler func(*Mesh, *broker.Message) error

func Wrap[I any](handler func(*Context, *I) error) MessageHandler {
	return func(m *Mesh, msg *broker.Message) error {

		envy := &envelope.Gate2MeshEnvelope{}
		if err := proto.Unmarshal(msg.Data, envy); err != nil {
			return err
		}

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
}
