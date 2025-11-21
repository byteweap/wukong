package logger

import "time"

// Entry TODO 日志条目接口, 用于链式调用
// 调用 Msg() Msgf()时执行打印
type Entry interface {
	Str(k string, v string) Entry
	Int64(k string, v int64) Entry
	Int(k string, v int) Entry
	Uint64(k string, v uint64) Entry
	Float(k string, v float64) Entry
	Bool(k string, v bool) Entry
	Time(k string, v time.Time) Entry
	Duration(k string, v time.Duration) Entry
	Any(k string, v any) Entry
	Err(err error) Entry
	Msg(message string)
	Msgf(format string, args ...any)
}
