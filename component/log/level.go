package log

import "strings"

// Level 表示日志级别。
type Level int8

// LevelKey 是日志级别字段名。
const LevelKey = "level"

const (
	// LevelDebug 表示 debug 级别。
	LevelDebug Level = iota - 1
	// LevelInfo 表示 info 级别。
	LevelInfo
	// LevelWarn 表示 warn 级别。
	LevelWarn
	// LevelError 表示 error 级别。
	LevelError
	// LevelFatal 表示 fatal 级别。
	LevelFatal
)

// Key 返回日志级别字段名。
func (l Level) Key() string {
	return LevelKey
}

// String 返回日志级别字符串。
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// ParseLevel 将字符串解析为日志级别。
func ParseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN":
		return LevelWarn
	case "ERROR":
		return LevelError
	case "FATAL":
		return LevelFatal
	}
	return LevelInfo
}
