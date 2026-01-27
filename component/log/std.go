package log

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

var _ Logger = (*stdLogger)(nil)

// stdLogger 对应标准库 [log.Logger] 的能力，并支持并发使用
type stdLogger struct {
	w         io.Writer
	isDiscard bool
	mu        sync.Mutex
	pool      *sync.Pool
}

// NewStdLogger 使用 writer 创建日志器
func NewStdLogger(w io.Writer) Logger {
	return &stdLogger{
		w:         w,
		isDiscard: w == io.Discard,
		pool: &sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

// Log 输出键值对日志
func (l *stdLogger) Log(level Level, keyvals ...any) error {
	if l.isDiscard || len(keyvals) == 0 {
		return nil
	}
	if (len(keyvals) & 1) == 1 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}

	buf := l.pool.Get().(*bytes.Buffer)
	defer l.pool.Put(buf)

	buf.WriteString(level.String())
	for i := 0; i < len(keyvals); i += 2 {
		_, _ = fmt.Fprintf(buf, " %s=%v", keyvals[i], keyvals[i+1])
	}
	buf.WriteByte('\n')
	defer buf.Reset()

	l.mu.Lock()
	defer l.mu.Unlock()
	_, err := l.w.Write(buf.Bytes())
	return err
}

// Close 兼容接口，当前无操作
func (l *stdLogger) Close() error {
	return nil
}
