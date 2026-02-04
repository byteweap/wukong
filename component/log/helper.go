package log

import (
	"context"
	"fmt"
	"os"
)

// DefaultMessageKey 默认消息字段名
var DefaultMessageKey = "msg"

// Option 是 Helper 的配置项
type Option func(*Helper)

// Helper 是日志辅助封装
type Helper struct {
	logger  Logger
	msgKey  string
	sprint  func(...any) string
	sprintf func(format string, a ...any) string
}

// WithMessageKey 设置消息字段名
func WithMessageKey(k string) Option {
	return func(opts *Helper) {
		opts.msgKey = k
	}
}

// WithSprint 设置 Sprint 实现
func WithSprint(sprint func(...any) string) Option {
	return func(opts *Helper) {
		opts.sprint = sprint
	}
}

// WithSprintf 设置 Sprintf 实现
func WithSprintf(sprintf func(format string, a ...any) string) Option {
	return func(opts *Helper) {
		opts.sprintf = sprintf
	}
}

// NewHelper 创建日志辅助器
func NewHelper(logger Logger, opts ...Option) *Helper {
	options := &Helper{
		msgKey:  DefaultMessageKey, // 默认消息字段名
		logger:  logger,
		sprint:  fmt.Sprint,
		sprintf: fmt.Sprintf,
	}
	for _, o := range opts {
		o(options)
	}
	return options
}

// WithContext 返回携带新 ctx 的浅拷贝，ctx 不能为空
func (h *Helper) WithContext(ctx context.Context) *Helper {
	return &Helper{
		msgKey:  h.msgKey,
		logger:  WithContext(ctx, h.logger),
		sprint:  h.sprint,
		sprintf: h.sprintf,
	}
}

// Enabled 判断级别是否可输出，基于底层 *Filter 进行判断
func (h *Helper) Enabled(level Level) bool {
	if l, ok := h.logger.(*Filter); ok {
		return level >= l.level
	}
	return true
}

// Logger 返回内部 logger
func (h *Helper) Logger() Logger {
	return h.logger
}

// Log 按级别输出键值对
func (h *Helper) Log(level Level, kvs ...any) {
	_ = h.logger.Log(level, kvs...)
}

// Debug 输出 debug 级别日志
func (h *Helper) Debug(a ...any) {
	if !h.Enabled(LevelDebug) {
		return
	}
	_ = h.logger.Log(LevelDebug, h.msgKey, h.sprint(a...))
}

// Debugf 输出 debug 级别格式化日志
func (h *Helper) Debugf(format string, a ...any) {
	if !h.Enabled(LevelDebug) {
		return
	}
	_ = h.logger.Log(LevelDebug, h.msgKey, h.sprintf(format, a...))
}

// Debugw 输出 debug 级别键值对日志
func (h *Helper) Debugw(kvs ...any) {
	_ = h.logger.Log(LevelDebug, kvs...)
}

// Info 输出 info 级别日志
func (h *Helper) Info(a ...any) {
	if !h.Enabled(LevelInfo) {
		return
	}
	_ = h.logger.Log(LevelInfo, h.msgKey, h.sprint(a...))
}

// Infof 输出 info 级别格式化日志
func (h *Helper) Infof(format string, a ...any) {
	if !h.Enabled(LevelInfo) {
		return
	}
	_ = h.logger.Log(LevelInfo, h.msgKey, h.sprintf(format, a...))
}

// Infow 输出 info 级别键值对日志
func (h *Helper) Infow(kvs ...any) {
	_ = h.logger.Log(LevelInfo, kvs...)
}

// Warn 输出 warn 级别日志
func (h *Helper) Warn(a ...any) {
	if !h.Enabled(LevelWarn) {
		return
	}
	_ = h.logger.Log(LevelWarn, h.msgKey, h.sprint(a...))
}

// Warnf 输出 warn 级别格式化日志
func (h *Helper) Warnf(format string, a ...any) {
	if !h.Enabled(LevelWarn) {
		return
	}
	_ = h.logger.Log(LevelWarn, h.msgKey, h.sprintf(format, a...))
}

// Warnw 输出 warn 级别键值对日志
func (h *Helper) Warnw(kvs ...any) {
	_ = h.logger.Log(LevelWarn, kvs...)
}

// Error 输出 error 级别日志
func (h *Helper) Error(a ...any) {
	if !h.Enabled(LevelError) {
		return
	}
	_ = h.logger.Log(LevelError, h.msgKey, h.sprint(a...))
}

// Errorf 输出 error 级别格式化日志
func (h *Helper) Errorf(format string, a ...any) {
	if !h.Enabled(LevelError) {
		return
	}
	_ = h.logger.Log(LevelError, h.msgKey, h.sprintf(format, a...))
}

// Errorw 输出 error 级别键值对日志
func (h *Helper) Errorw(kvs ...any) {
	_ = h.logger.Log(LevelError, kvs...)
}

// Fatal 输出 fatal 级别日志并退出进程
func (h *Helper) Fatal(a ...any) {
	_ = h.logger.Log(LevelFatal, h.msgKey, h.sprint(a...))
	os.Exit(1)
}

// Fatalf 输出 fatal 级别格式化日志并退出进程
func (h *Helper) Fatalf(format string, a ...any) {
	_ = h.logger.Log(LevelFatal, h.msgKey, h.sprintf(format, a...))
	os.Exit(1)
}

// Fatalw 输出 fatal 级别键值对日志并退出进程
func (h *Helper) Fatalw(kvs ...any) {
	_ = h.logger.Log(LevelFatal, kvs...)
	os.Exit(1)
}
