package msgpack

import (
	"github.com/byteweap/wukong/encoding"
	"github.com/vmihailenco/msgpack/v5"
)

const Name = "msgpack"

var c = &codec{}

func init() {
	encoding.RegisterCodec(codec{})
}

type codec struct{}

// Name 编解码器名称
func (codec) Name() string {
	return Name
}

// Marshal 编码
func (codec) Marshal(v any) ([]byte, error) {
	return msgpack.Marshal(v)
}

// Unmarshal 解码
func (codec) Unmarshal(data []byte, v any) error {
	return msgpack.Unmarshal(data, v)
}

// Marshal 编码
func Marshal(v any) ([]byte, error) {
	return c.Marshal(v)
}

// Unmarshal 解码
func Unmarshal(data []byte, v any) error {
	return c.Unmarshal(data, v)
}
