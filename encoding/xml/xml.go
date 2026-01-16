package xml

import (
	"encoding/xml"

	"github.com/byteweap/wukong/encoding"
)

// Name is the name registered for the xml codec.
const Name = "xml"

var c = &codec{}

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with xml.
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
