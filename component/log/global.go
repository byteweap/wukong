package log

import (
	"context"
	"fmt"
	"os"
	"sync"
)

// globalLogger 作为进程内的全局日志器。
var global = &loggerAppliance{}

// loggerAppliance 作为 Logger 的代理，更新后会影响所有子 logger。
type loggerAppliance struct {
	lock sync.RWMutex
	Logger
}

func init() {
	global.SetLogger(DefaultLogger)
}

func (a *loggerAppliance) SetLogger(in Logger) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Logger = in
}

// SetLogger 建议在任何日志调用前设置，且非线程安全。
func SetLogger(logger Logger) {
	global.SetLogger(logger)
}

// GetLogger 返回当前进程的全局日志器。
func GetLogger() Logger {
	global.lock.RLock()
	defer global.lock.RUnlock()
	return global.Logger
}

// Log 按级别输出键值对日志。
func Log(level Level, keyvals ...any) {
	_ = global.Log(level, keyvals...)
}

// Context 返回带上下文的日志辅助器。
func Context(ctx context.Context) *Helper {
	return NewHelper(WithContext(ctx, global.Logger))
}

// Debug 输出 debug 级别日志。
func Debug(a ...any) {
	_ = global.Log(LevelDebug, DefaultMessageKey, fmt.Sprint(a...))
}

// Debugf 输出 debug 级别格式化日志。
func Debugf(format string, a ...any) {
	_ = global.Log(LevelDebug, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Debugw 输出 debug 级别键值对日志。
func Debugw(keyvals ...any) {
	_ = global.Log(LevelDebug, keyvals...)
}

// Info 输出 info 级别日志。
func Info(a ...any) {
	_ = global.Log(LevelInfo, DefaultMessageKey, fmt.Sprint(a...))
}

// Infof 输出 info 级别格式化日志。
func Infof(format string, a ...any) {
	_ = global.Log(LevelInfo, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Infow 输出 info 级别键值对日志。
func Infow(keyvals ...any) {
	_ = global.Log(LevelInfo, keyvals...)
}

// Warn 输出 warn 级别日志。
func Warn(a ...any) {
	_ = global.Log(LevelWarn, DefaultMessageKey, fmt.Sprint(a...))
}

// Warnf 输出 warn 级别格式化日志。
func Warnf(format string, a ...any) {
	_ = global.Log(LevelWarn, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Warnw 输出 warn 级别键值对日志。
func Warnw(keyvals ...any) {
	_ = global.Log(LevelWarn, keyvals...)
}

// Error 输出 error 级别日志。
func Error(a ...any) {
	_ = global.Log(LevelError, DefaultMessageKey, fmt.Sprint(a...))
}

// Errorf 输出 error 级别格式化日志。
func Errorf(format string, a ...any) {
	_ = global.Log(LevelError, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Errorw 输出 error 级别键值对日志。
func Errorw(keyvals ...any) {
	_ = global.Log(LevelError, keyvals...)
}

// Fatal 输出 fatal 级别日志并退出进程。
func Fatal(a ...any) {
	_ = global.Log(LevelFatal, DefaultMessageKey, fmt.Sprint(a...))
	os.Exit(1)
}

// Fatalf 输出 fatal 级别格式化日志并退出进程。
func Fatalf(format string, a ...any) {
	_ = global.Log(LevelFatal, DefaultMessageKey, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Fatalw 输出 fatal 级别键值对日志并退出进程。
func Fatalw(keyvals ...any) {
	_ = global.Log(LevelFatal, keyvals...)
	os.Exit(1)
}
