package logger

// Logger 日志接口
type Logger interface {
	// With 添加字段，返回新的 Logger
	With(string, string) Logger
	// Debug 创建 Debug 级别日志条目
	Debug() Entry
	// Info 创建 Info 级别日志条目
	Info() Entry
	// Warn 创建 Warn 级别日志条目
	Warn() Entry
	// Error 创建 Error 级别日志条目
	Error() Entry
	// Fatal 创建 Fatal 级别日志条目
	Fatal() Entry
	// Panic 创建 Panic 级别日志条目
	Panic() Entry
}
