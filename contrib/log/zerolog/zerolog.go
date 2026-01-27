package zerolog

import (
	"github.com/byteweap/wukong/component/log"
	"github.com/rs/zerolog"
)

// Logger 日志记录器
type Logger struct {
	log *zerolog.Logger
}

var _ log.Logger = (*Logger)(nil)

// NewLogger 创建日志记录器
func NewLogger(logger *zerolog.Logger) log.Logger {
	return &Logger{
		log: logger,
	}
}

// Log 发送日志
func (l *Logger) Log(level log.Level, kvs ...any) (err error) {

	var event *zerolog.Event

	if len(kvs) == 0 {
		return nil
	}
	if len(kvs)%2 != 0 {
		kvs = append(kvs, "")
	}

	switch level {
	case log.LevelDebug:
		event = l.log.Debug()
	case log.LevelInfo:
		event = l.log.Info()
	case log.LevelWarn:
		event = l.log.Warn()
	case log.LevelError:
		event = l.log.Error()
	case log.LevelFatal:
		event = l.log.Fatal()
	default:
		event = l.log.Debug()
	}

	for i := 0; i < len(kvs); i += 2 {
		key, ok := kvs[i].(string)
		if !ok {
			continue
		}
		event = event.Any(key, kvs[i+1])
	}
	event.Send()
	return
}
