package kconv

import (
	"bytes"
	"encoding/binary"
	"reflect"

	"github.com/byteweap/wukong/pkg/kcodec/json"
)

func Byte(data any) byte {
	return Uint8(data)
}

func Bytes(data any) []byte {
	if data == nil {
		return nil
	}

	var (
		err error
		buf = bytes.NewBuffer(nil)
	)

	switch v := data.(type) {
	case int:
		err = binary.Write(buf, binary.BigEndian, int64(v))
	case *int:
		err = binary.Write(buf, binary.BigEndian, int64(*v))
	case uint:
		err = binary.Write(buf, binary.BigEndian, uint64(v))
	case *uint:
		err = binary.Write(buf, binary.BigEndian, uint64(*v))
	case bool, *bool, int8, *int8, int16, *int16, int32, *int32, int64, *int64, uint8, *uint8, uint16, *uint16, uint32, *uint32, uint64, *uint64, float32, *float32, float64, *float64:
		err = binary.Write(buf, binary.BigEndian, v)
	case uintptr:
		err = binary.Write(buf, binary.BigEndian, uint64(v))
	case *uintptr:
		err = binary.Write(buf, binary.BigEndian, uint64(*v))
	case complex64, *complex64, complex128, *complex128:
		return nil
	case string:
		return StringToBytes(v)
	case *string:
		return StringToBytes(*v)
	case []byte:
		return v
	case *[]byte:
		return *v
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
		case reflect.Invalid:
			return nil
		case reflect.Bool:
			err = binary.Write(buf, binary.BigEndian, rv.Bool())
		case reflect.String:
			return StringToBytes(rv.String())
		case reflect.Int, reflect.Int64:
			err = binary.Write(buf, binary.BigEndian, rv.Int())
		case reflect.Int8:
			err = binary.Write(buf, binary.BigEndian, int8(rv.Int()))
		case reflect.Int16:
			err = binary.Write(buf, binary.BigEndian, int16(rv.Int()))
		case reflect.Int32:
			err = binary.Write(buf, binary.BigEndian, int32(rv.Int()))
		case reflect.Uint, reflect.Uint64, reflect.Uintptr:
			err = binary.Write(buf, binary.BigEndian, rv.Uint())
		case reflect.Uint8:
			err = binary.Write(buf, binary.BigEndian, uint8(rv.Uint()))
		case reflect.Uint16:
			err = binary.Write(buf, binary.BigEndian, uint16(rv.Uint()))
		case reflect.Uint32:
			err = binary.Write(buf, binary.BigEndian, uint32(rv.Uint()))
		case reflect.Float32, reflect.Float64:
			err = binary.Write(buf, binary.BigEndian, rv.Float())
		case reflect.Complex64, reflect.Complex128:
			return nil
		default:
			b, err := json.Marshal(data)
			if err != nil {
				return nil
			}
			return b
		}
	}
	if err != nil {
		return nil
	}

	return buf.Bytes()
}

func BytePointer(data any) *byte {
	v := Byte(data)
	return &v
}

func BytesPointer(data any) *[]byte {
	v := Bytes(data)
	return &v
}
