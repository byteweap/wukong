package logrus

import (
	"github.com/sirupsen/logrus"

	"github.com/byteweap/wukong/component/log"
)

// Logger 日志记录器
type Logger struct {
	log *logrus.Logger
}

var _ log.Logger = (*Logger)(nil)

// NewLogger 创建日志记录器
func NewLogger(logger *logrus.Logger) log.Logger {
	return &Logger{
		log: logger,
	}
}

// Log 发送日志
func (l *Logger) Log(level log.Level, kvs ...any) (err error) {
	var (
		logrusLevel logrus.Level
		fields      logrus.Fields = make(map[string]any)
		msg         string
	)

	switch level {
	case log.LevelDebug:
		logrusLevel = logrus.DebugLevel
	case log.LevelInfo:
		logrusLevel = logrus.InfoLevel
	case log.LevelWarn:
		logrusLevel = logrus.WarnLevel
	case log.LevelError:
		logrusLevel = logrus.ErrorLevel
	case log.LevelFatal:
		logrusLevel = logrus.FatalLevel
	default:
		logrusLevel = logrus.DebugLevel
	}

	if logrusLevel > l.log.Level {
		return
	}

	if len(kvs) == 0 {
		return nil
	}
	if len(kvs)%2 != 0 {
		kvs = append(kvs, "")
	}
	for i := 0; i < len(kvs); i += 2 {
		key, ok := kvs[i].(string)
		if !ok {
			continue
		}
		if key == logrus.FieldKeyMsg {
			msg, _ = kvs[i+1].(string)
			continue
		}
		fields[key] = kvs[i+1]
	}

	if len(fields) > 0 {
		l.log.WithFields(fields).Log(logrusLevel, msg)
	} else {
		l.log.Log(logrusLevel, msg)
	}

	return
}
