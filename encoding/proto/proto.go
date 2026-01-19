// Package proto 定义 protobuf 编解码器。导入此包将自动注册编解码器。
package proto

import (
	"errors"
	"reflect"

	"github.com/byteweap/wukong/encoding"
	"google.golang.org/protobuf/proto"
)

// Name 是 proto 编解码器注册的名称。
const Name = "proto"

var c = &codec{}

func init() {
	encoding.RegisterCodec(codec{})
}

// codec 是基于 protobuf 的编解码器实现，是传输的默认编解码器。
type codec struct{}

func (codec) Marshal(v any) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (codec) Unmarshal(data []byte, v any) error {
	pm, err := getProtoMessage(v)
	if err != nil {
		return err
	}
	return proto.Unmarshal(data, pm)
}

func (codec) Name() string {
	return Name
}

func getProtoMessage(v any) (proto.Message, error) {
	if msg, ok := v.(proto.Message); ok {
		return msg, nil
	}
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return nil, errors.New("not proto message")
	}

	val = val.Elem()
	return getProtoMessage(val.Interface())
}

func Marshal(v any) ([]byte, error) {
	return c.Marshal(v)
}

func Unmarshal(data []byte, v any) error {
	return c.Unmarshal(data, v)
}
