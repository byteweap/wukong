package zerolog

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/byteweap/wukong/component/logger"
)

// Log 日志实现
type Log struct {
	log zerolog.Logger
}

// 确保 Log 实现 logger.Logger 接口
var _ logger.Logger = (*Log)(nil)

// New 创建新的日志实例
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

// With 添加字段，返回新的 Logger
func (l *Log) With(k, v string) logger.Logger {
	return &Log{
		log: l.log.With().Str(k, v).Logger(),
	}
}

// Debug 创建 Debug 级别日志条目
func (l *Log) Debug() logger.Entry {
	return newEntry(l.log.Debug())
}

// Info 创建 Info 级别日志条目
func (l *Log) Info() logger.Entry {
	return newEntry(l.log.Info())
}

// Warn 创建 Warn 级别日志条目
func (l *Log) Warn() logger.Entry {
	return newEntry(l.log.Warn())
}

// Error 创建 Error 级别日志条目
func (l *Log) Error() logger.Entry {
	return newEntry(l.log.Error())
}

// Fatal 创建 Fatal 级别日志条目
func (l *Log) Fatal() logger.Entry {
	return newEntry(l.log.Fatal())
}

// Panic 创建 Panic 级别日志条目
func (l *Log) Panic() logger.Entry {
	return newEntry(l.log.Panic())
}

// newOutput 根据配置创建输出器
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

// convertLevel 将字符串级别转换为 zerolog.Level
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
