package logger

import "time"

// Entry 日志条目接口，用于链式调用
type Entry interface {
	// Str 添加字符串字段
	Str(k string, v string) Entry
	// Int64 添加 int64 字段
	Int64(k string, v int64) Entry
	// Int 添加 int 字段
	Int(k string, v int) Entry
	// Uint64 添加 uint64 字段
	Uint64(k string, v uint64) Entry
	// Float 添加 float64 字段
	Float(k string, v float64) Entry
	// Bool 添加 bool 字段
	Bool(k string, v bool) Entry
	// Time 添加时间字段
	Time(k string, v time.Time) Entry
	// Duration 添加时长字段
	Duration(k string, v time.Duration) Entry
	// Any 添加任意类型字段
	Any(k string, v any) Entry
	// Err 添加错误字段
	Err(err error) Entry
	// Msg 打印日志消息
	Msg(message string)
	// Msgf 格式化打印日志消息
	Msgf(format string, args ...any)
}
