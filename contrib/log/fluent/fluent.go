package fluent

import (
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/fluent/fluent-logger-golang/fluent"

	"github.com/byteweap/wukong/component/log"
)

// Logger 是 fluent 日志记录器
type Logger struct {
	opts options
	log  *fluent.Fluent
}

var _ log.Logger = (*Logger)(nil)

// NewLogger 创建 fluent 日志记录器
// target
//
//	tcp://127.0.0.1:24224
//	unix://var/run/fluent/fluent.sock
func NewLogger(addr string, opts ...Option) (*Logger, error) {
	option := options{}
	for _, o := range opts {
		o(&option)
	}
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	c := fluent.Config{
		Timeout:            option.timeout,
		WriteTimeout:       option.writeTimeout,
		BufferLimit:        option.bufferLimit,
		RetryWait:          option.retryWait,
		MaxRetry:           option.maxRetry,
		MaxRetryWait:       option.maxRetryWait,
		TagPrefix:          option.tagPrefix,
		Async:              option.async,
		ForceStopAsyncSend: option.forceStopAsyncSend,
	}
	switch u.Scheme {
	case "tcp":
		host, port, err2 := net.SplitHostPort(u.Host)
		if err2 != nil {
			return nil, err2
		}
		if c.FluentPort, err = strconv.Atoi(port); err != nil {
			return nil, err
		}
		c.FluentNetwork = u.Scheme
		c.FluentHost = host
	case "unix":
		c.FluentNetwork = u.Scheme
		c.FluentSocketPath = u.Path
	default:
		return nil, fmt.Errorf("unknown network: %s", u.Scheme)
	}
	fl, err := fluent.New(c)
	if err != nil {
		return nil, err
	}
	return &Logger{
		opts: option,
		log:  fl,
	}, nil
}

// Log 发送键值对日志
func (l *Logger) Log(level log.Level, kvs ...any) error {
	if len(kvs) == 0 {
		return nil
	}
	if len(kvs)%2 != 0 {
		kvs = append(kvs, "KEYVALS UNPAIRED")
	}

	data := make(map[string]string, len(kvs)/2+1)

	for i := 0; i < len(kvs); i += 2 {
		data[fmt.Sprint(kvs[i])] = fmt.Sprint(kvs[i+1])
	}

	if err := l.log.Post(level.String(), data); err != nil {
		println(err)
	}
	return nil
}

// Close 关闭日志记录器
func (l *Logger) Close() error {
	return l.log.Close()
}
