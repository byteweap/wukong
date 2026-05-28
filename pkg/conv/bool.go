package conv

import (
	"reflect"
	"strings"
	"time"
	"unsafe"
)

func Bool(data any) bool {
	if data == nil {
		return false
	}

	switch v := data.(type) {
	case []byte:
		return stringToBool(*(*string)(unsafe.Pointer(&v)))
	case *[]byte:
		return v != nil && stringToBool(*(*string)(unsafe.Pointer(v)))
	case time.Time:
		return v.IsZero()
	case *time.Time:
		return v != nil && v.IsZero()
	default:
		return reflectBool(indirectValue(data))
	}
}

func stringToBool(v string) bool {
	return v != "" && v != "0" && strings.ToLower(v) != "false"
}

func reflectBool(rv reflect.Value) bool {
	if !rv.IsValid() {
		return false
	}

	switch rv.Kind() {
	case reflect.Bool:
		return rv.Bool()
	case reflect.String:
		return stringToBool(rv.String())
	case reflect.Uintptr:
		return rv.Uint() != 0
	case reflect.UnsafePointer:
		return !rv.IsNil() && uint(rv.Pointer()) != 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() != 0
	case reflect.Complex64, reflect.Complex128:
		return stringToBool(String(rv.Complex()))
	case reflect.Array:
		return rv.Len() != 0
	case reflect.Slice, reflect.Map:
		return !rv.IsNil() && rv.Len() != 0
	case reflect.Struct:
		return true
	case reflect.Chan, reflect.Func, reflect.Interface:
		return !rv.IsNil()
	default:
		return false
	}
}

func Bools(data any) []bool {
	return convertSlice(data, Bool)
}

func BoolPointer(any any) *bool {
	v := Bool(any)
	return &v
}

func BoolsPointer(any any) *[]bool {
	v := Bools(any)
	return &v
}
