package xml

import (
	"encoding/xml"

	"github.com/byteweap/wukong/encoding"
)

// Name 是 XML 编解码器注册的名称。
const Name = "xml"

var c = &codec{}

func init() {
	encoding.RegisterCodec(codec{})
}

// codec 是基于 XML 的编解码器实现。
type codec struct{}

func (codec) Marshal(v any) ([]byte, error) {
	return xml.Marshal(v)
}

func (codec) Unmarshal(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}

func Marshal(v any) ([]byte, error) {
	return c.Marshal(v)
}

func Unmarshal(data []byte, v any) error {
	return c.Unmarshal(data, v)
}
