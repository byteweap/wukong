package yaml

import (
	"github.com/byteweap/wukong/encoding"
	"gopkg.in/yaml.v3"
)

// Name is the name registered for the yaml codec.
const Name = "yaml"

var c = &codec{}

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with yaml.
type codec struct{}

func (codec) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (codec) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
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
