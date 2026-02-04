package log

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	// DefaultCaller 返回调用方的文件与行号
	DefaultCaller = Caller(4)

	// DefaultTimestamp 返回当前时间戳
	DefaultTimestamp = Timestamp(time.RFC3339)
)

// Valuer 返回日志值。
type Valuer func(ctx context.Context) any

// Value 返回 Valuer 的结果或原值
func Value(ctx context.Context, v any) any {
	if v, ok := v.(Valuer); ok {
		return v(ctx)
	}
	return v
}

// Caller 返回调用方的 pkg/file:line 描述
func Caller(depth int) Valuer {
	return func(context.Context) any {
		_, file, line, _ := runtime.Caller(depth)
		idx := strings.LastIndexByte(file, '/')
		if idx == -1 {
			return file[idx+1:] + ":" + strconv.Itoa(line)
		}
		idx = strings.LastIndexByte(file[:idx], '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

// Timestamp 返回指定格式的时间戳 Valuer
func Timestamp(layout string) Valuer {
	return func(context.Context) any {
		return time.Now().Format(layout)
	}
}

// bindValues 将 Valuer 绑定为实际值
func bindValues(ctx context.Context, kvs []any) {
	for i := 1; i < len(kvs); i += 2 {
		if v, ok := kvs[i].(Valuer); ok {
			kvs[i] = v(ctx)
		}
	}
}

// containsValuer 判断键值对中是否包含 Valuer
func containsValuer(kvs []any) bool {
	for i := 1; i < len(kvs); i += 2 {
		if _, ok := kvs[i].(Valuer); ok {
			return true
		}
	}
	return false
}
