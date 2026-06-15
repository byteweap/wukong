package log

import "strings"

// Level 表示日志级别
type Level int8

// LevelKey 是日志级别字段名
const LevelKey = "level"

const (
	levelDebugString = "DEBUG"
	levelInfoString  = "INFO"
	levelWarnString  = "WARN"
	levelErrorString = "ERROR"
	levelFatalString = "FATAL"
)

const (
	// LevelDebug 表示 debug 级别
	LevelDebug Level = iota - 1
	// LevelInfo 表示 info 级别
	LevelInfo
	// LevelWarn 表示 warn 级别
	LevelWarn
	// LevelError 表示 error 级别
	LevelError
	// LevelFatal 表示 fatal 级别
	LevelFatal
)

// Key 返回日志级别字段名
func (l Level) Key() string {
	return LevelKey
}

// String 返回日志级别字符串
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return levelDebugString
	case LevelInfo:
		return levelInfoString
	case LevelWarn:
		return levelWarnString
	case LevelError:
		return levelErrorString
	case LevelFatal:
		return levelFatalString
	default:
		return ""
	}
}

// ParseLevel 将字符串解析为日志级别
func ParseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case levelDebugString:
		return LevelDebug
	case levelInfoString:
		return LevelInfo
	case levelWarnString:
		return LevelWarn
	case levelErrorString:
		return LevelError
	case levelFatalString:
		return LevelFatal
	}
	return LevelInfo
}
