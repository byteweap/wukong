package kconv

import "reflect"

func Uint16(data any) uint16 {
	return uint16(Uint64(data))
}

func Uint16s(data any) (slice []uint16) {
	if data == nil {
		return
	}

	switch v := data.(type) {
	case []int:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]int:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case []int8:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]int8:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case []int16:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]int16:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case []int32:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]int32:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case []int64:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]int64:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case []uint:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]uint:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case []uint8:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]uint8:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case []uint16:
		return v
	case *[]uint16:
		return *v
	case []uint32:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]uint32:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case []uint64:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]uint64:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case []float32:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]float32:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case []float64:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]float64:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case []complex64:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]complex64:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case []complex128:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]complex128:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case []string:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]string:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case []bool:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]bool:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case []any:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[]any:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
	case [][]byte:
		slice = make([]uint16, len(v))
		for i := range v {
			slice[i] = Uint16(v[i])
		}
	case *[][]byte:
		slice = make([]uint16, len(*v))
		for i := range *v {
			slice[i] = Uint16((*v)[i])
		}
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
		case reflect.Slice, reflect.Array:
			count := rv.Len()
			slice = make([]uint16, count)
			for i := 0; i < count; i++ {
				slice[i] = Uint16(rv.Index(i).Interface())
			}
		}
	}

	return
}

func Uint16Pointer(data any) *uint16 {
	v := Uint16(data)
	return &v
}

func Uint16sPointer(data any) *[]uint16 {
	v := Uint16s(data)
	return &v
}
