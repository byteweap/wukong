package zerolog

import "time"

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
	LevelPanic = "panic"
)

const (
	ModeConsole = "console"
	ModeFile    = "file"
)

const (
	// default value
	defaultLevel            = LevelDebug
	defaultLevelFieldName   = "level"
	defaultTimeFieldName    = "time"
	defaultMessageFieldName = "message"
	defaultTimeFormat       = time.RFC3339
	defaultMode             = ModeConsole

	// default value about file
	defaultFilename   = "./logs.log"
	defaultMaxSize    = 100 // MB
	defaultMaxBackups = 10
	defaultMaxAge     = 30
	defaultCompress   = false
	defaultLocalTime  = false
)

type fileOptions struct {
	filename   string // 文件名 默认: ./logs.log
	maxSize    int    // 单文件最大容量 单位:MB 默认: 10MB
	maxBackups int    // 最大备份文件数 单位:个 默认: 10天
	maxAge     int    // 最大保留时间 单位:天 默认: 10天
	compress   bool   // 备份是否压缩 默认: false
	localTime  bool   // 使用本地时区时间, 默认: false (UTC)
}

type options struct {
	level            string      // 日志等级, 默认: debug, 可选: debug,info,warn,error
	levelFieldName   string      // 日志等级字段名, 默认: level
	timeFieldName    string      // 时间字段名, 默认: time
	messageFieldName string      // 日志内容字段名, 默认: message
	timeFormat       string      // 时间格式, 默认: 2006-01-02T15:04:05Z07:00 (time.RFC3339)
	mode             string      // 日志模式, 默认: console, 可选: console,file
	fileOpts         fileOptions // 日志模式为file时,相关配置 默认: ./logs.log, 10MB, 10个文件, 10天, 不压缩, UTC时间
}

type Option func(*options)

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

// WithLevel 指定日志等级, 默认debug, 可选: debug,info,warn,error
func WithLevel(level string) Option {
	return func(o *options) {
		if level == LevelDebug ||
			level == LevelInfo ||
			level == LevelWarn ||
			level == LevelError {
			o.level = level
		}
	}
}

// WithMode 指定输出模式, 默认console, 可选: console, file
func WithMode(mode string) Option {
	return func(o *options) {
		if mode == ModeConsole || mode == ModeFile {
			o.mode = mode
		}
	}
}

// WithLevelFieldName 指定日志等级字段名, 默认: level
func WithLevelFieldName(v string) Option {
	return func(o *options) {
		if v != "" {
			o.levelFieldName = v
		}
	}
}

// WithTimeFieldName 指定时间字段名, 默认: time
func WithTimeFieldName(v string) Option {
	return func(o *options) {
		if v != "" {
			o.timeFieldName = v
		}
	}
}

// WithTimeFormat 指定时间格式, 默认: 2006-01-02T15:04:05Z07:00 (time.RFC3339)
func WithTimeFormat(v string) Option {
	return func(o *options) {
		if v != "" {
			o.timeFormat = v
		}
	}
}

// WithMessageFieldName 指定日志内容字段名, 默认: message
func WithMessageFieldName(v string) Option {
	return func(o *options) {
		if v != "" {
			o.messageFieldName = v
		}
	}
}

// WithFile 指定日志文件相关配置, 默认: ./logs.log, 10MB, 10个文件, 10天, 不压缩, UTC时间
func WithFile(filename string, maxSize, maxBackups, maxAge int, compress, localTime bool) Option {
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
