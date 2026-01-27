package log

import "io"

type writerWrapper struct {
	helper *Helper
	level  Level
}

// WriterOptionFn 是 writerWrapper 的配置项
type WriterOptionFn func(w *writerWrapper)

// WithWriterLevel 设置 writerWrapper 的日志级别
func WithWriterLevel(level Level) WriterOptionFn {
	return func(w *writerWrapper) {
		w.level = level
	}
}

// WithWriteMessageKey 设置 writerWrapper 的消息字段名
func WithWriteMessageKey(key string) WriterOptionFn {
	return func(w *writerWrapper) {
		w.helper.msgKey = key
	}
}

// NewWriter 返回一个 writer 包装器
func NewWriter(logger Logger, opts ...WriterOptionFn) io.Writer {
	ww := &writerWrapper{
		helper: NewHelper(logger, WithMessageKey(DefaultMessageKey)),
		level:  LevelInfo, // 默认级别
	}
	for _, opt := range opts {
		opt(ww)
	}
	return ww
}

func (ww *writerWrapper) Write(p []byte) (int, error) {
	ww.helper.Log(ww.level, ww.helper.msgKey, string(p))
	return 0, nil
}
