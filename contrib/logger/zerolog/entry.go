package zerolog

import (
	"time"

	"github.com/rs/zerolog"

	"github.com/byteweap/wukong/component/logger"
)

// Entry 日志条目实现
type Entry struct {
	e *zerolog.Event
}

// 确保 Entry 实现 logger.Entry 接口
var _ logger.Entry = (*Entry)(nil)

// newEntry 创建新的日志条目
func newEntry(e *zerolog.Event) *Entry {
	return &Entry{e: e}
}

// Str 添加字符串字段
func (e *Entry) Str(k, v string) logger.Entry {
	e.e.Str(k, v)
	return e
}

// Int64 添加 int64 字段
func (e *Entry) Int64(k string, v int64) logger.Entry {
	e.e.Int64(k, v)
	return e
}

// Int 添加 int 字段
func (e *Entry) Int(k string, v int) logger.Entry {
	e.e.Int(k, v)
	return e
}

// Uint64 添加 uint64 字段
func (e *Entry) Uint64(k string, v uint64) logger.Entry {
	e.e.Uint64(k, v)
	return e
}

// Float 添加 float64 字段
func (e *Entry) Float(k string, v float64) logger.Entry {
	e.e.Float64(k, v)
	return e
}

// Bool 添加 bool 字段
func (e *Entry) Bool(k string, v bool) logger.Entry {
	e.e.Bool(k, v)
	return e
}

// Time 添加时间字段
func (e *Entry) Time(k string, v time.Time) logger.Entry {
	e.e.Time(k, v)
	return e
}

// Duration 添加时长字段
func (e *Entry) Duration(k string, v time.Duration) logger.Entry {
	e.e.Dur(k, v)
	return e
}

// Any 添加任意类型字段
func (e *Entry) Any(k string, v any) logger.Entry {
	e.e.Any(k, v)
	return e
}

// Err 添加错误字段
func (e *Entry) Err(err error) logger.Entry {
	e.e.Err(err)
	return e
}

// Msg 打印日志消息
func (e *Entry) Msg(message string) {
	e.e.Msg(message)
}

// Msgf 格式化打印日志消息
func (e *Entry) Msgf(format string, args ...any) {
	e.e.Msgf(format, args...)
}
