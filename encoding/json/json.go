package json

import (
	"encoding/json"
	"reflect"

	"github.com/byteweap/wukong/encoding"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Name 是 JSON 编解码器注册的名称。
const Name = "json"

var (
	c = &codec{}

	// MarshalOptions 可配置的 JSON 序列化选项。
	MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: true,
	}
	// UnmarshalOptions 可配置的 JSON 反序列化选项。
	UnmarshalOptions = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

func init() {
	encoding.RegisterCodec(codec{})
}

// codec 是基于 JSON 的编解码器实现。
type codec struct{}

func (codec) Name() string {
	return Name
}

func (codec) Marshal(v any) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	case proto.Message:
		return MarshalOptions.Marshal(m)
	default:
		return json.Marshal(m)
	}
}

func (codec) Unmarshal(data []byte, v any) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	case proto.Message:
		return UnmarshalOptions.Unmarshal(data, m)
	default:
		rv := reflect.ValueOf(v)
		for rv := rv; rv.Kind() == reflect.Ptr; {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		if m, ok := reflect.Indirect(rv).Interface().(proto.Message); ok {
			return UnmarshalOptions.Unmarshal(data, m)
		}
		return json.Unmarshal(data, m)
	}
}

func Marshal(v any) ([]byte, error) {
	return c.Marshal(v)
}

func Unmarshal(data []byte, v any) error {
	return c.Unmarshal(data, v)
}
