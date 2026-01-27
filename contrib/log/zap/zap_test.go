package zap

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/byteweap/wukong/component/log"
)

// testWriteSyncer 测试写入器
type testWriteSyncer struct {
	output []string
}

// Write 写入数据
func (x *testWriteSyncer) Write(p []byte) (n int, err error) {
	x.output = append(x.output, string(p))
	return len(p), nil
}

// Sync 同步数据
func (x *testWriteSyncer) Sync() error {
	return nil
}

// TestLogger 测试日志记录器
func TestLogger(t *testing.T) {
	syncer := &testWriteSyncer{}
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), syncer, zap.DebugLevel)
	zlogger := zap.New(core).WithOptions()
	logger := NewLogger(zlogger)

	defer func() { _ = logger.Close() }()

	zlog := log.NewHelper(logger)

	zlog.Debugw("log", "debug")
	zlog.Infow("log", "info")
	zlog.Warnw("log", "warn")
	zlog.Errorw("log", "error")
	zlog.Errorw("log", "error", "except warn")
	zlog.Info("hello world")

	except := []string{
		"{\"level\":\"debug\",\"msg\":\"\",\"log\":\"debug\"}\n",
		"{\"level\":\"info\",\"msg\":\"\",\"log\":\"info\"}\n",
		"{\"level\":\"warn\",\"msg\":\"\",\"log\":\"warn\"}\n",
		"{\"level\":\"error\",\"msg\":\"\",\"log\":\"error\"}\n",
		"{\"level\":\"warn\",\"msg\":\"Keyvalues must appear in pairs: [log error except warn]\"}\n",
		"{\"level\":\"info\",\"msg\":\"hello world\"}\n", // not {"level":"info","msg":"","msg":"hello world"}
	}
	for i, s := range except {
		if s != syncer.output[i] {
			t.Logf("except=%s, got=%s", s, syncer.output[i])
			t.Fail()
		}
	}
}
