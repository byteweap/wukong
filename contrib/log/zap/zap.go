package zap

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/byteweap/wukong/component/log"
)

// Option 配置函数
type Option func(*Logger)

// WithMessageKey 设置消息键
func WithMessageKey(key string) Option {
	return func(l *Logger) {
		l.msgKey = key
	}
}

// Logger 日志记录器
type Logger struct {
	log    *zap.Logger
	msgKey string
}

var _ log.Logger = (*Logger)(nil)

// NewLogger 创建日志记录器
func NewLogger(zlog *zap.Logger) *Logger {
	return &Logger{
		log:    zlog,
		msgKey: log.DefaultMessageKey,
	}
}

// Log 发送日志
func (l *Logger) Log(level log.Level, kvs ...any) error {
	// 若该级别被禁用则跳过格式化开销
	if zapcore.Level(level) < zapcore.DPanicLevel && !l.log.Core().Enabled(zapcore.Level(level)) {
		return nil
	}
	var (
		msg    = ""
		keylen = len(kvs)
	)
	if keylen == 0 || keylen%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", kvs))
		return nil
	}

	data := make([]zap.Field, 0, (keylen/2)+1)
	for i := 0; i < keylen; i += 2 {
		if kvs[i].(string) == l.msgKey {
			msg, _ = kvs[i+1].(string)
			continue
		}
		data = append(data, zap.Any(fmt.Sprint(kvs[i]), kvs[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug(msg, data...)
	case log.LevelInfo:
		l.log.Info(msg, data...)
	case log.LevelWarn:
		l.log.Warn(msg, data...)
	case log.LevelError:
		l.log.Error(msg, data...)
	case log.LevelFatal:
		l.log.Fatal(msg, data...)
	}
	return nil
}

// Sync 刷新日志
func (l *Logger) Sync() error {
	return l.log.Sync()
}

// Close 关闭日志记录器
func (l *Logger) Close() error {
	return l.Sync()
}
