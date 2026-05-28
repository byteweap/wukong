package conv

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strconv"
	"time"
)

func Int64(data any) int64 {
	if data == nil {
		return 0
	}

	switch v := data.(type) {
	case int:
		return int64(v)
	case *int:
		return int64(*v)
	case int8:
		return int64(v)
	case *int8:
		return int64(*v)
	case int16:
		return int64(v)
	case *int16:
		return int64(*v)
	case int32:
		return int64(v)
	case *int32:
		return int64(*v)
	case int64:
		return v
	case *int64:
		return *v
	case uint:
		return int64(v)
	case *uint:
		return int64(*v)
	case uint8:
		return int64(v)
	case *uint8:
		return int64(*v)
	case uint16:
		return int64(v)
	case *uint16:
		return int64(*v)
	case uint32:
		return int64(v)
	case *uint32:
		return int64(*v)
	case uint64:
		return int64(v)
	case *uint64:
		return int64(*v)
	case float32:
		return int64(v)
	case *float32:
		return int64(*v)
	case float64:
		return int64(v)
	case *float64:
		return int64(*v)
	case complex64:
		return int64(real(v))
	case *complex64:
		return int64(real(*v))
	case complex128:
		return int64(real(v))
	case *complex128:
		return int64(real(*v))
	case bool:
		if v {
			return 1
		}
		return 0
	case *bool:
		if *v {
			return 1
		}
		return 0
	case time.Time:
		return v.UnixNano()
	case *time.Time:
		return v.UnixNano()
	case []byte:
		buf := make([]byte, 8)
		copy(buf[len(buf)-len(v):], v)

		var i int64
		if err := binary.Read(bytes.NewReader(buf), binary.BigEndian, &i); err == nil {
			return i
		} else {
			return 0
		}
	case *[]byte:
		return Int64(*v)
	default:
		var (
			rv   = reflect.ValueOf(data)
			kind = rv.Kind()
		)

		for kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}

		switch kind {
		case reflect.Bool:
			return Int64(rv.Bool())
		case reflect.String:
			i, _ := strconv.ParseInt(rv.String(), 10, 64)
			return i
		case reflect.Uintptr:
			return int64(rv.Uint())
		case reflect.UnsafePointer:
			return int64(rv.Pointer())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return rv.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return int64(rv.Uint())
		case reflect.Float32, reflect.Float64:
			return int64(rv.Float())
		case reflect.Complex64, reflect.Complex128:
			return int64(real(rv.Complex()))
		default:
			return 0
		}
	}
}

func Int64s(data any) []int64 {
	return convertSlice(data, Int64)
}

func Int64Pointer(any any) *int64 {
	v := Int64(any)
	return &v
}

func Int64sPointer(any any) *[]int64 {
	v := Int64s(any)
	return &v
}
