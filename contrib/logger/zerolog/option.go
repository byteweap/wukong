package zerolog

import "time"

// 日志级别常量
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
	LevelPanic = "panic"
)

// 日志输出模式常量
const (
	ModeConsole = "console"
	ModeFile    = "file"
)

const (
	// 默认值
	defaultLevel            = LevelDebug
	defaultLevelFieldName   = "level"
	defaultTimeFieldName    = "time"
	defaultMessageFieldName = "message"
	defaultTimeFormat       = time.RFC3339
	defaultMode             = ModeConsole

	// 文件相关默认值
	defaultFilename   = "./logs.log"
	defaultMaxSize    = 100 // MB
	defaultMaxBackups = 10
	defaultMaxAge     = 30
	defaultCompress   = false
	defaultLocalTime  = false
)

// fileOptions 文件输出配置
type fileOptions struct {
	filename   string // 文件名，默认: ./logs.log
	maxSize    int    // 单文件最大容量，单位:MB，默认: 100MB
	maxBackups int    // 最大备份文件数，默认: 10
	maxAge     int    // 最大保留时间，单位:天，默认: 30天
	compress   bool   // 备份是否压缩，默认: false
	localTime  bool   // 使用本地时区时间，默认: false (UTC)
}

// options 日志配置选项
type options struct {
	level            string      // 日志等级，默认: debug，可选: debug,info,warn,error
	levelFieldName   string      // 日志等级字段名，默认: level
	timeFieldName    string      // 时间字段名，默认: time
	messageFieldName string      // 日志内容字段名，默认: message
	timeFormat       string      // 时间格式，默认: 2006-01-02T15:04:05Z07:00 (time.RFC3339)
	mode             string      // 日志模式，默认: console，可选: console,file
	fileOpts         fileOptions // 日志模式为file时的相关配置
}

// Option 配置选项函数类型
type Option func(*options)

// defaultOptions 返回默认配置
func defaultOptions() *options {
	return &options{
		level:            defaultLevel,
		levelFieldName:   defaultLevelFieldName,
		timeFieldName:    defaultTimeFieldName,
		messageFieldName: defaultMessageFieldName,
		timeFormat:       defaultTimeFormat,
		mode:             defaultMode,
		fileOpts: fileOptions{
			filename:   defaultFilename,
			maxSize:    defaultMaxSize,
			maxBackups: defaultMaxBackups,
			maxAge:     defaultMaxAge,
			compress:   defaultCompress,
			localTime:  defaultLocalTime,
		},
	}
}

// Level 指定日志等级，默认: debug，可选: debug,info,warn,error
func Level(level string) Option {
	return func(o *options) {
		if level == LevelDebug ||
			level == LevelInfo ||
			level == LevelWarn ||
			level == LevelError {
			o.level = level
		}
	}
}

// Mode 指定输出模式，默认: console，可选: console,file
func Mode(mode string) Option {
	return func(o *options) {
		if mode == ModeConsole || mode == ModeFile {
			o.mode = mode
		}
	}
}

// LevelFieldName 指定日志等级字段名，默认: level
func LevelFieldName(v string) Option {
	return func(o *options) {
		if v != "" {
			o.levelFieldName = v
		}
	}
}

// TimeFieldName 指定时间字段名，默认: time
func TimeFieldName(v string) Option {
	return func(o *options) {
		if v != "" {
			o.timeFieldName = v
		}
	}
}

// TimeFormat 指定时间格式，默认: 2006-01-02T15:04:05Z07:00 (time.RFC3339)
func TimeFormat(v string) Option {
	return func(o *options) {
		if v != "" {
			o.timeFormat = v
		}
	}
}

// MessageFieldName 指定日志内容字段名，默认: message
func MessageFieldName(v string) Option {
	return func(o *options) {
		if v != "" {
			o.messageFieldName = v
		}
	}
}

// File 指定日志文件相关配置
func File(filename string, maxSize, maxBackups, maxAge int, compress, localTime bool) Option {
	return func(o *options) {
		if filename != "" {
			o.fileOpts.filename = filename
		}
		if maxSize > 0 {
			o.fileOpts.maxSize = maxSize
		}
		if maxBackups > 0 {
			o.fileOpts.maxBackups = maxBackups
		}
		if maxAge > 0 {
			o.fileOpts.maxAge = maxAge
		}
		o.fileOpts.compress = compress
		o.fileOpts.localTime = localTime
	}
}
