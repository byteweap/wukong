package encoding

import (
	"strings"
)

// Codec 定义编解码接口，用于序列化和反序列化消息。
// 实现必须线程安全，方法可能被并发调用。
type Codec interface {
	// Name 返回编解码器名称，用作传输时的内容类型标识。
	// 返回值必须静态不变。
	Name() string
	// Marshal 将 v 序列化为字节数组。
	Marshal(v any) ([]byte, error)
	// Unmarshal 将字节数组反序列化到 v。
	Unmarshal(data []byte, v any) error
}

var registeredCodecs = make(map[string]Codec)

// RegisterCodec 注册编解码器，供所有传输客户端和服务器使用。
func RegisterCodec(codec Codec) {
	if codec == nil {
		panic("cannot register a nil Codec")
	}
	if codec.Name() == "" {
		panic("cannot register Codec with empty string result for Name()")
	}
	contentSubtype := strings.ToLower(codec.Name())
	registeredCodecs[contentSubtype] = codec
}

// GetCodec 根据内容子类型获取已注册的编解码器，未注册则返回 nil。
//
// content-subtype 应为小写。
func GetCodec(contentSubtype string) Codec {
	return registeredCodecs[contentSubtype]
}
