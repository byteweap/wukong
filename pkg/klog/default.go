package klog

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type defaultLog struct {
	log zerolog.Logger
}

var _ Logger = (*defaultLog)(nil)

func newDefaultLog(opts ...Option) *defaultLog {

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

	return &defaultLog{
		log: zerolog.New(output).With().Timestamp().CallerWithSkipFrameCount(3).Logger(),
	}
}
func (d *defaultLog) With(k string, v string) Logger {
	return &defaultLog{
		log: d.log.With().Str(k, v).Logger(),
	}
}

func (d *defaultLog) Debug() Entry {
	return newDefaultEntry(d.log.Debug())
}

func (d *defaultLog) Info() Entry {
	return newDefaultEntry(d.log.Info())
}

func (d *defaultLog) Warn() Entry {
	return newDefaultEntry(d.log.Warn())
}

func (d *defaultLog) Error() Entry {
	return newDefaultEntry(d.log.Error())
}

func (d *defaultLog) Fatal() Entry {
	return newDefaultEntry(d.log.Fatal())
}

func (d *defaultLog) Panic() Entry {
	return newDefaultEntry(d.log.Panic())
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
		return zerolog.NoLevel
	}
}

type defaultEntry struct {
	e *zerolog.Event
}

var _ Entry = (*defaultEntry)(nil)

func newDefaultEntry(e *zerolog.Event) *defaultEntry {
	return &defaultEntry{
		e: e,
	}
}

func (d *defaultEntry) Str(k string, v string) Entry {
	d.e.Str(k, v)
	return d
}

func (d *defaultEntry) Int64(k string, v int64) Entry {
	d.e.Int64(k, v)
	return d
}

func (d *defaultEntry) Int(k string, v int) Entry {
	d.e.Int(k, v)
	return d
}

func (d *defaultEntry) Uint64(k string, v uint64) Entry {
	d.e.Uint64(k, v)
	return d
}

func (d *defaultEntry) Float(k string, v float64) Entry {
	d.e.Float64(k, v)
	return d
}

func (d *defaultEntry) Bool(k string, v bool) Entry {
	d.e.Bool(k, v)
	return d
}

func (d *defaultEntry) Time(k string, v time.Time) Entry {
	d.e.Time(k, v)
	return d
}

func (d *defaultEntry) Duration(k string, v time.Duration) Entry {
	d.e.Dur(k, v)
	return d
}

func (d *defaultEntry) Any(k string, v any) Entry {
	d.e.Any(k, v)
	return d
}

func (d *defaultEntry) Err(err error) Entry {
	d.e.Err(err)
	return d
}

func (d *defaultEntry) Msg(message string) {
	d.e.Msg(message)
}

func (d *defaultEntry) Msgf(format string, args ...any) {
	d.e.Msgf(format, args...)
}
