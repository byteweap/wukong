package kcodec

import (
	"errors"

	"github.com/byteweap/wukong/pkg/kcodec/json"
	"github.com/byteweap/wukong/pkg/kcodec/msgpack"
	"github.com/byteweap/wukong/pkg/kcodec/proto"
	"github.com/byteweap/wukong/pkg/kcodec/toml"
	"github.com/byteweap/wukong/pkg/kcodec/xml"
	"github.com/byteweap/wukong/pkg/kcodec/yml"
)

var codecs = make(map[string]Codec)

func init() {
	_ = Register(json.Codec)
	_ = Register(proto.Codec)
	_ = Register(toml.Codec)
	_ = Register(xml.Codec)
	_ = Register(yml.Codec)
	_ = Register(msgpack.Codec)
}

type Codec interface {
	// Name 编解码器类型
	Name() string
	// Marshal 编码
	Marshal(v any) ([]byte, error)
	// Unmarshal 解码
	Unmarshal(data []byte, v any) error
}

// Register 注册编解码器
func Register(codec Codec) error {
	if codec == nil {
		return errors.New("invalid codec")
	}
	name := codec.Name()
	if name == "" {
		return errors.New("invalid codec name")
	}
	codecs[name] = codec
	return nil
}

// Invoke 调用编解码器
func Invoke(name string) (Codec, error) {
	codec, ok := codecs[name]
	if !ok {
		return nil, errors.New("codec is not registered")
	}
	return codec, nil
}
