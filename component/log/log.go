package log

import (
	"context"
	"log"
)

// DefaultLogger 是默认日志器。
var DefaultLogger = NewStdLogger(log.Writer())

// Logger 是日志接口。
type Logger interface {
	Log(level Level, keyvals ...any) error
}

type logger struct {
	logger    Logger
	prefix    []any
	hasValuer bool
	ctx       context.Context
}

func (c *logger) Log(level Level, keyvals ...any) error {
	kvs := make([]any, 0, len(c.prefix)+len(keyvals))
	kvs = append(kvs, c.prefix...)
	if c.hasValuer {
		bindValues(c.ctx, kvs)
	}
	kvs = append(kvs, keyvals...)
	return c.logger.Log(level, kvs...)
}

// With 追加日志字段。
func With(l Logger, kv ...any) Logger {
	c, ok := l.(*logger)
	if !ok {
		return &logger{logger: l, prefix: kv, hasValuer: containsValuer(kv), ctx: context.Background()}
	}
	kvs := make([]any, 0, len(c.prefix)+len(kv))
	kvs = append(kvs, c.prefix...)
	kvs = append(kvs, kv...)
	return &logger{
		logger:    c.logger,
		prefix:    kvs,
		hasValuer: containsValuer(kvs),
		ctx:       c.ctx,
	}
}

// WithContext 返回携带新 ctx 的浅拷贝，ctx 不能为空。
func WithContext(ctx context.Context, l Logger) Logger {
	switch v := l.(type) {
	default:
		return &logger{logger: l, ctx: ctx}
	case *logger:
		lv := *v
		lv.ctx = ctx
		return &lv
	case *Filter:
		fv := *v
		fv.logger = WithContext(ctx, fv.logger)
		return &fv
	}
}
