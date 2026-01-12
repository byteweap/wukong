package zerolog

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/byteweap/wukong/component/logger"
)

type Log struct {
	log zerolog.Logger
}

var _ logger.Logger = (*Log)(nil)

func New(opts ...Option) *Log {

	ops := defaultOptions()
	for _, opt := range opts {
		opt(ops)
	}

	output := newOutput(ops)

	zerolog.SetGlobalLevel(convertLevel(ops.level))
	zerolog.TimeFieldFormat = ops.timeFormat
	zerolog.TimestampFieldName = ops.timeFieldName
	zerolog.LevelFieldName = ops.levelFieldName
	zerolog.MessageFieldName = ops.messageFieldName

	return &Log{
		log: zerolog.New(output).With().Timestamp().CallerWithSkipFrameCount(3).Logger(),
	}
}

func (l *Log) With(k, v string) logger.Logger {
	return &Log{
		log: l.log.With().Str(k, v).Logger(),
	}
}

func (l *Log) Debug() logger.Entry {
	return newEntry(l.log.Debug())
}

func (l *Log) Info() logger.Entry {
	return newEntry(l.log.Info())
}

func (l *Log) Warn() logger.Entry {
	return newEntry(l.log.Warn())
}

func (l *Log) Error() logger.Entry {
	return newEntry(l.log.Error())
}

func (l *Log) Fatal() logger.Entry {
	return newEntry(l.log.Fatal())
}

func (l *Log) Panic() logger.Entry {
	return newEntry(l.log.Panic())
}

func newOutput(ops *options) io.Writer {
	switch ops.mode {
	case ModeConsole:
		return zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: ops.timeFormat,
		}
	case ModeFile:
		return &lumberjack.Logger{
			Filename:   ops.fileOpts.filename,
			MaxSize:    ops.fileOpts.maxSize,
			MaxAge:     ops.fileOpts.maxAge,
			MaxBackups: ops.fileOpts.maxBackups,
			LocalTime:  ops.fileOpts.localTime,
			Compress:   ops.fileOpts.compress,
		}
	}
	return os.Stdout
}

func convertLevel(level string) zerolog.Level {
	switch level {
	case LevelDebug:
		return zerolog.DebugLevel
	case LevelInfo:
		return zerolog.InfoLevel
	case LevelWarn:
		return zerolog.WarnLevel
	case LevelError:
		return zerolog.ErrorLevel
	case LevelFatal:
		return zerolog.FatalLevel
	case LevelPanic:
		return zerolog.PanicLevel
	default:
		return zerolog.DebugLevel
	}
}
