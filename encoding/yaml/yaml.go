package yaml

import (
	"github.com/byteweap/wukong/encoding"
	"gopkg.in/yaml.v3"
)

// Name 是 YAML 编解码器注册的名称。
const Name = "yaml"

var c = &codec{}

func init() {
	encoding.RegisterCodec(codec{})
}

// codec 是基于 YAML 的编解码器实现。
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
