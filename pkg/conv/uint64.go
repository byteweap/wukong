package conv

import (
	"reflect"
	"strconv"
	"time"
)

func Uint64(data any) uint64 {
	if data == nil {
		return 0
	}

	switch v := data.(type) {
	case int:
		return uint64(v)
	case *int:
		return uint64(*v)
	case int8:
		return uint64(v)
	case *int8:
		return uint64(*v)
	case int16:
		return uint64(v)
	case *int16:
		return uint64(*v)
	case int32:
		return uint64(v)
	case *int32:
		return uint64(*v)
	case int64:
		return uint64(v)
	case *int64:
		return uint64(*v)
	case uint:
		return uint64(v)
	case *uint:
		return uint64(*v)
	case uint8:
		return uint64(v)
	case *uint8:
		return uint64(*v)
	case uint16:
		return uint64(v)
	case *uint16:
		return uint64(*v)
	case uint32:
		return uint64(v)
	case *uint32:
		return uint64(*v)
	case uint64:
		return v
	case *uint64:
		return *v
	case float32:
		return uint64(v)
	case *float32:
		return uint64(*v)
	case float64:
		return uint64(v)
	case *float64:
		return uint64(*v)
	case complex64:
		return uint64(real(v))
	case *complex64:
		return uint64(real(*v))
	case complex128:
		return uint64(real(v))
	case *complex128:
		return uint64(real(*v))
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
		return uint64(v.UnixNano())
	case *time.Time:
		return uint64(v.UnixNano())
	default:
		var (
			rv   = reflect.ValueOf(data)
			kind = rv.Kind()
		)

		for kind == reflect.Pointer {
			rv = rv.Elem()
			kind = rv.Kind()
		}

		switch kind {
		case reflect.Bool:
			return Uint64(rv.Bool())
		case reflect.String:
			i, _ := strconv.ParseUint(rv.String(), 0, 64)
			return i
		case reflect.Uintptr:
			return rv.Uint()
		case reflect.UnsafePointer:
			return uint64(rv.Pointer())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Uint64(rv.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return rv.Uint()
		case reflect.Float32, reflect.Float64:
			return uint64(rv.Float())
		case reflect.Complex64, reflect.Complex128:
			return uint64(real(rv.Complex()))
		default:
			return 0
		}
	}
}

func Uint64s(data any) []uint64 {
	return convertSlice(data, Uint64)
}

func Uint64Pointer(data any) *uint64 {
	v := Uint64(data)
	return &v
}

func Uint64sPointer(data any) *[]uint64 {
	v := Uint64s(data)
	return &v
}
