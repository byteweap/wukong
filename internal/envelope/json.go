package envelope

type JsonEnvelope struct {
	Seq     uint64 `json:"seq"`     // 客户端消息唯一标识,客户端递增
	Cmd     uint32 `json:"cmd"`     // 指令
	Payload []byte `json:"payload"` // 有效载荷
}

var _ Envelope = (*JsonEnvelope)(nil)

func (j JsonEnvelope) GetSeq() uint64 {
	return j.Seq
}

func (j JsonEnvelope) GetCmd() uint32 {
	return j.Cmd
}

func (j JsonEnvelope) GetPayload() []byte {
	return j.Payload
}
