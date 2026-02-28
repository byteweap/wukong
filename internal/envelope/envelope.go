package envelope

type Envelope interface {
	GetSeq() uint64
	GetCmd() uint32
	GetPayload() []byte
}
