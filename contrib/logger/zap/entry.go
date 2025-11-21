package logrus

import (
	"time"

	"github.com/byteweap/wukong/plugin/logger"
)

// Entry TODO
type Entry struct {
}

var _ logger.Entry = (*Entry)(nil)

func NewEntry() *Entry {
	return &Entry{}
}

func (e *Entry) Str(k string, v string) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *Entry) Int64(k string, v int64) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *Entry) Int(k string, v int) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *Entry) Uint64(k string, v uint64) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *Entry) Float(k string, v float64) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *Entry) Bool(k string, v bool) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *Entry) Time(k string, v time.Time) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *Entry) Duration(k string, v time.Duration) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *Entry) Any(k string, v any) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *Entry) Err(err error) logger.Entry {
	//TODO implement me
	panic("implement me")
}

func (e *Entry) Msg(message string) {
	//TODO implement me
	panic("implement me")
}

func (e *Entry) Msgf(format string, args ...any) {
	//TODO implement me
	panic("implement me")
}
