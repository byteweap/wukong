package logger

type Logger interface {
	// With 添加字段,返回新的Logger
	With(string, string) Logger
	Debug() Entry
	Info() Entry
	Warn() Entry
	Error() Entry
	Fatal() Entry
	Panic() Entry
}
