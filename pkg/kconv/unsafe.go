package kconv

import (
	"unsafe"
)

// StringToBytes 字符串无拷贝转字节数组
func StringToBytes(s string) []byte {
	data := unsafe.StringData(s)
	return unsafe.Slice(data, len(s))
}

// BytesToString 字节数组无拷贝转字符串
func BytesToString(b []byte) string {
	return unsafe.String(&b[0], len(b))
}
