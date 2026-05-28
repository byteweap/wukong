package log

import (
	"testing"
)

func TestInfo(_ *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp)
	logger = With(logger, "caller", DefaultCaller)
	_ = logger.Log(LevelInfo, "key1", "value1")
}

func TestStdLogger_Close(t *testing.T) {
	stdLog := NewStdLogger(nil)
	// stdLogger 实现了 Close 方法
	if closer, ok := stdLog.(interface{ Close() error }); ok {
		err := closer.Close()
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	}
}

func TestHelper_WithSprint(t *testing.T) {
	customSprint := func(a ...any) string {
		return "custom"
	}
	log := NewHelper(DefaultLogger, WithSprint(customSprint))
	log.Debug("test")
}

func TestHelper_WithSprintf(t *testing.T) {
	customSprintf := func(format string, a ...any) string {
		return "custom"
	}
	log := NewHelper(DefaultLogger, WithSprintf(customSprintf))
	log.Debugf("test %s", "arg")
}

func TestHelper_Logger(t *testing.T) {
	log := NewHelper(DefaultLogger)
	if log.Logger() == nil {
		t.Error("expected non-nil logger")
	}
}

func TestHelper_Log(t *testing.T) {
	log := NewHelper(DefaultLogger)
	log.Log(LevelInfo, "key", "value")
}

func TestHelper_Infow(t *testing.T) {
	log := NewHelper(DefaultLogger)
	log.Infow("key", "value")
}

func TestHelper_Warnw(t *testing.T) {
	log := NewHelper(DefaultLogger)
	log.Warnw("key", "value")
}

func TestHelper_Error(t *testing.T) {
	log := NewHelper(DefaultLogger)
	log.Error("test error")
}

func TestHelper_Errorw(t *testing.T) {
	log := NewHelper(DefaultLogger)
	log.Errorw("key", "value")
}

func TestHelper_Enabled(t *testing.T) {
	filter := NewFilter(DefaultLogger, FilterLevel(LevelInfo))
	log := NewHelper(filter)

	if log.Enabled(LevelDebug) {
		t.Error("expected LevelDebug to be disabled")
	}
	if !log.Enabled(LevelInfo) {
		t.Error("expected LevelInfo to be enabled")
	}
	if !log.Enabled(LevelError) {
		t.Error("expected LevelError to be enabled")
	}

	// Test with non-filter logger
	log2 := NewHelper(DefaultLogger)
	if !log2.Enabled(LevelDebug) {
		t.Error("expected LevelDebug to be enabled for non-filter logger")
	}
}

func TestHelper_Info_Disabled(t *testing.T) {
	filter := NewFilter(DefaultLogger, FilterLevel(LevelWarn))
	log := NewHelper(filter)
	// Should not log because level is below threshold
	log.Info("should not be logged")
	log.Infof("should not be logged %s", "test")
}

func TestHelper_Warn_Disabled(t *testing.T) {
	filter := NewFilter(DefaultLogger, FilterLevel(LevelError))
	log := NewHelper(filter)
	// Should not log because level is below threshold
	log.Warn("should not be logged")
	log.Warnf("should not be logged %s", "test")
}

func TestHelper_Error_Disabled(t *testing.T) {
	filter := NewFilter(DefaultLogger, FilterLevel(LevelFatal))
	log := NewHelper(filter)
	// Should not log because level is below threshold
	log.Error("should not be logged")
	log.Errorf("should not be logged %s", "test")
}
