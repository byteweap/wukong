package kconv

import "github.com/byteweap/wukong/pkg/kcodec/json"

// Scan json反序列化
func Scan(data any, target any) error {
	return json.Unmarshal(Bytes(data), target)
}
