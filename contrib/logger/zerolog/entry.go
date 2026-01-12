package zerolog

import (
	"time"

	"github.com/rs/zerolog"

	"github.com/byteweap/wukong/component/logger"
)

type Entry struct {
	e *zerolog.Event
}

var _ logger.Entry = (*Entry)(nil)

func newEntry(e *zerolog.Event) *Entry {
	return &Entry{e: e}
}

func (e *Entry) Str(k, v string) logger.Entry {
	e.e.Str(k, v)
	return e
}

func (e *Entry) Int64(k string, v int64) logger.Entry {
	e.e.Int64(k, v)
	return e
}

func (e *Entry) Int(k string, v int) logger.Entry {
	e.e.Int(k, v)
	return e
}

func (e *Entry) Uint64(k string, v uint64) logger.Entry {
	e.e.Uint64(k, v)
	return e
}

func (e *Entry) Float(k string, v float64) logger.Entry {
	e.e.Float64(k, v)
	return e
}

func (e *Entry) Bool(k string, v bool) logger.Entry {
	e.e.Bool(k, v)
	return e
}

func (e *Entry) Time(k string, v time.Time) logger.Entry {
	e.e.Time(k, v)
	return e
}

func (e *Entry) Duration(k string, v time.Duration) logger.Entry {
	e.e.Dur(k, v)
	return e
}

func (e *Entry) Any(k string, v any) logger.Entry {
	e.e.Any(k, v)
	return e
}

func (e *Entry) Err(err error) logger.Entry {
	e.e.Err(err)
	return e
}

func (e *Entry) Msg(message string) {
	e.e.Msg(message)
}

func (e *Entry) Msgf(format string, args ...any) {
	e.e.Msgf(format, args...)
}
