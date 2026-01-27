package fluent

import (
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/byteweap/wukong/component/log"
)

// TestMain 启动测试监听器
func TestMain(m *testing.M) {
	listener := func(ln net.Listener) {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_, err = io.ReadAll(conn)
		if err != nil {
			return
		}
	}

	if ln, err := net.Listen("tcp", ":24224"); err == nil {
		defer ln.Close()
		go func() {
			for {
				listener(ln)
			}
		}()
	}

	os.Exit(m.Run())
}

// TestWithTimeout 测试设置超时
func TestWithTimeout(t *testing.T) {
	opts := new(options)
	var duration time.Duration = 1000000000
	funcTimeout := WithTimeout(duration)
	funcTimeout(opts)
	if opts.timeout != duration {
		t.Errorf("WithTimeout() = %v, want %v", opts.timeout, duration)
	}
}

// TestWithWriteTimeout 测试设置写入超时
func TestWithWriteTimeout(t *testing.T) {
	opts := new(options)
	var duration time.Duration = 1000000000
	funcWriteTimeout := WithWriteTimeout(duration)
	funcWriteTimeout(opts)
	if opts.writeTimeout != duration {
		t.Errorf("WithWriteTimeout() = %v, want %v", opts.writeTimeout, duration)
	}
}

// TestWithBufferLimit 测试设置缓冲上限
func TestWithBufferLimit(t *testing.T) {
	opts := new(options)
	bufferLimit := 1000000000
	funcBufferLimit := WithBufferLimit(bufferLimit)
	funcBufferLimit(opts)
	if opts.bufferLimit != bufferLimit {
		t.Errorf("WithBufferLimit() = %d, want %d", opts.bufferLimit, bufferLimit)
	}
}

// TestWithRetryWait 测试设置重试间隔
func TestWithRetryWait(t *testing.T) {
	opts := new(options)
	retryWait := 1000000000
	funcRetryWait := WithRetryWait(retryWait)
	funcRetryWait(opts)
	if opts.retryWait != retryWait {
		t.Errorf("WithRetryWait() = %d, want %d", opts.retryWait, retryWait)
	}
}

// TestWithMaxRetry 测试设置最大重试次数
func TestWithMaxRetry(t *testing.T) {
	opts := new(options)
	maxRetry := 1000000000
	funcMaxRetry := WithMaxRetry(maxRetry)
	funcMaxRetry(opts)
	if opts.maxRetry != maxRetry {
		t.Errorf("WithMaxRetry() = %d, want %d", opts.maxRetry, maxRetry)
	}
}

// TestWithMaxRetryWait 测试设置最大重试等待
func TestWithMaxRetryWait(t *testing.T) {
	opts := new(options)
	maxRetryWait := 1000000000
	funcMaxRetryWait := WithMaxRetryWait(maxRetryWait)
	funcMaxRetryWait(opts)
	if opts.maxRetryWait != maxRetryWait {
		t.Errorf("WithMaxRetryWait() = %d, want %d", opts.maxRetryWait, maxRetryWait)
	}
}

// TestWithTagPrefix 测试设置标签前缀
func TestWithTagPrefix(t *testing.T) {
	opts := new(options)
	tagPrefix := "tag_prefix"
	funcTagPrefix := WithTagPrefix(tagPrefix)
	funcTagPrefix(opts)
	if opts.tagPrefix != tagPrefix {
		t.Errorf("WithTagPrefix() = %s, want %s", opts.tagPrefix, tagPrefix)
	}
}

// TestWithAsync 测试设置异步发送
func TestWithAsync(t *testing.T) {
	opts := new(options)
	async := true
	funcAsync := WithAsync(async)
	funcAsync(opts)
	if opts.async != async {
		t.Errorf("WithAsync() = %t, want %t", opts.async, async)
	}
}

// TestWithForceStopAsyncSend 测试强制停止异步发送
func TestWithForceStopAsyncSend(t *testing.T) {
	opts := new(options)
	forceStopAsyncSend := true
	funcForceStopAsyncSend := WithForceStopAsyncSend(forceStopAsyncSend)
	funcForceStopAsyncSend(opts)
	if opts.forceStopAsyncSend != forceStopAsyncSend {
		t.Errorf("WithForceStopAsyncSend() = %t, want %t", opts.forceStopAsyncSend, forceStopAsyncSend)
	}
}

// TestLogger 测试日志记录器
func TestLogger(t *testing.T) {
	logger, err := NewLogger("tcp://127.0.0.1:24224")
	if err != nil {
		t.Error(err)
	}
	defer logger.Close()
	flog := log.NewHelper(logger)

	flog.Debug("log", "test")
	flog.Info("log", "test")
	flog.Warn("log", "test")
	flog.Error("log", "test")
}

// TestLoggerWithOpt 测试带配置的日志记录器
func TestLoggerWithOpt(t *testing.T) {
	var duration time.Duration = 1000000000
	logger, err := NewLogger("tcp://127.0.0.1:24224", WithTimeout(duration))
	if err != nil {
		t.Error(err)
	}
	defer logger.Close()
	flog := log.NewHelper(logger)

	flog.Debug("log", "test")
	flog.Info("log", "test")
	flog.Warn("log", "test")
	flog.Error("log", "test")
}

// TestLoggerError 测试错误地址
func TestLoggerError(t *testing.T) {
	errCase := []string{
		"foo",
		"tcp://127.0.0.1/",
		"tcp://127.0.0.1:1234a",
		"tcp://127.0.0.1:65535",
		"https://127.0.0.1:8080",
		"unix://foo/bar",
	}
	for _, errc := range errCase {
		_, err := NewLogger(errc)
		if err == nil {
			t.Error(err)
		}
	}
}

// BenchmarkLoggerPrint 基准测试日志输出
func BenchmarkLoggerPrint(b *testing.B) {
	b.SetParallelism(100)
	logger, err := NewLogger("tcp://127.0.0.1:24224")
	flog := log.NewHelper(logger)
	if err != nil {
		b.Error(err)
	}
	defer logger.Close()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			flog.Info("log", "test")
		}
	})
}

// BenchmarkLoggerHelperV 基准测试日志助手
func BenchmarkLoggerHelperV(b *testing.B) {
	b.SetParallelism(100)
	logger, err := NewLogger("tcp://127.0.0.1:24224")
	if err != nil {
		b.Error(err)
	}
	h := log.NewHelper(logger)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Info("log", "test")
		}
	})
}

// BenchmarkLoggerHelperInfo 基准测试信息日志
func BenchmarkLoggerHelperInfo(b *testing.B) {
	b.SetParallelism(100)
	logger, err := NewLogger("tcp://127.0.0.1:24224")
	if err != nil {
		b.Error(err)
	}
	defer logger.Close()
	h := log.NewHelper(logger)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Info("test")
		}
	})
}

// BenchmarkLoggerHelperInfof 基准测试格式化日志
func BenchmarkLoggerHelperInfof(b *testing.B) {
	b.SetParallelism(100)
	logger, err := NewLogger("tcp://127.0.0.1:24224")
	if err != nil {
		b.Error(err)
	}
	defer logger.Close()
	h := log.NewHelper(logger)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Infof("log %s", "test")
		}
	})
}

// BenchmarkLoggerHelperInfow 基准测试键值日志
func BenchmarkLoggerHelperInfow(b *testing.B) {
	b.SetParallelism(100)
	logger, err := NewLogger("tcp://127.0.0.1:24224")
	if err != nil {
		b.Error(err)
	}
	defer logger.Close()
	h := log.NewHelper(logger)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Infow("log", "test")
		}
	})
}
