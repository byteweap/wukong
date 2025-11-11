package klog

import "time"

// Logger 日志接口,支持链式调用
// 默认基于zerolog实现, 可实现Logger、Entry自行扩展
type Logger interface {
	With(string, string) Logger
	Debug() Entry
	Info() Entry
	Warn() Entry
	Error() Entry
	Fatal() Entry
	Panic() Entry
}

// Entry 日志条目接口, 用于链式调用
// 只有在执行 Msg() / Msgf() 才会打印
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

func New(opts ...Option) Logger {
	return newDefaultLog(opts...)
}
