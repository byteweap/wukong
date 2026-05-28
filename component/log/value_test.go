package log

import (
	"context"
	"testing"
)

func TestValue(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	_ = logger.Log(LevelInfo, "msg", "helloworld")

	logger = DefaultLogger
	logger = With(logger)
	_ = logger.Log(LevelDebug, "msg", "helloworld")

	var v1 any
	got := Value(context.Background(), v1)
	if got != v1 {
		t.Errorf("Value() = %v, want %v", got, v1)
	}
	var v2 Valuer = func(context.Context) any {
		return 3
	}
	got = Value(context.Background(), v2)
	res := got.(int)
	if res != 3 {
		t.Errorf("Value() = %v, want %v", res, 3)
	}
}

func TestTimestamp(t *testing.T) {
	ts := Timestamp("2006-01-02")
	val := ts(context.Background())
	if val == nil {
		t.Error("expected non-nil timestamp")
	}
	if _, ok := val.(string); !ok {
		t.Error("expected string timestamp")
	}
}

func TestCaller(t *testing.T) {
	c := Caller(1)
	val := c(context.Background())
	if val == nil {
		t.Error("expected non-nil caller")
	}
	if _, ok := val.(string); !ok {
		t.Error("expected string caller")
	}
}

func TestBindValues(t *testing.T) {
	kvs := []any{"key1", "value1", "key2", DefaultTimestamp}
	bindValues(context.Background(), kvs)
	if kvs[3] == nil {
		t.Error("expected timestamp to be bound")
	}
}

func TestContainsValuer(t *testing.T) {
	kvs1 := []any{"key1", "value1"}
	if containsValuer(kvs1) {
		t.Error("expected false for non-valuer kvs")
	}

	kvs2 := []any{"key1", DefaultTimestamp}
	if !containsValuer(kvs2) {
		t.Error("expected true for valuer kvs")
	}
}
