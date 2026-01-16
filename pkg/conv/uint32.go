package conv

import "reflect"

func Uint32(data any) uint32 {
	return uint32(Uint64(data))
}

func Uint32s(data any) (slice []uint32) {
	if data == nil {
		return
	}

	switch v := data.(type) {
	case []int:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]int:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case []int8:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]int8:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case []int16:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]int16:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case []int32:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]int32:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case []int64:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]int64:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case []uint:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]uint:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case []uint8:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]uint8:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case []uint16:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]uint16:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case []uint32:
		return v
	case *[]uint32:
		return *v
	case []uint64:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]uint64:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case []float32:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]float32:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case []float64:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]float64:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case []complex64:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]complex64:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case []complex128:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]complex128:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case []string:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]string:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case []bool:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]bool:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case []any:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[]any:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
		}
	case [][]byte:
		slice = make([]uint32, len(v))
		for i := range v {
			slice[i] = Uint32(v[i])
		}
	case *[][]byte:
		slice = make([]uint32, len(*v))
		for i := range *v {
			slice[i] = Uint32((*v)[i])
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
			slice = make([]uint32, count)
			for i := 0; i < count; i++ {
				slice[i] = Uint32(rv.Index(i).Interface())
			}
		}
	}

	return
}

func Uint32Pointer(data any) *uint32 {
	v := Uint32(data)
	return &v
}

func Uint32sPointer(data any) *[]uint32 {
	v := Uint32s(data)
	return &v
}
