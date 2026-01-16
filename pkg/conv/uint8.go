package conv

import "reflect"

func Uint8(data any) uint8 {
	return uint8(Uint64(data))
}

func Uint8s(data any) (slice []uint8) {
	if data == nil {
		return
	}

	switch v := data.(type) {
	case []int:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]int:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case []int8:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]int8:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case []int16:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]int16:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case []int32:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]int32:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case []int64:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]int64:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case []uint:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]uint:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case []uint8:
		return v
	case *[]uint8:
		return *v
	case []uint16:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]uint16:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case []uint32:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]uint32:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case []uint64:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]uint64:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case []float32:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]float32:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case []float64:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]float64:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case []complex64:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]complex64:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case []complex128:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]complex128:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case []string:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]string:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case []bool:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]bool:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case []any:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[]any:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
		}
	case [][]byte:
		slice = make([]uint8, len(v))
		for i := range v {
			slice[i] = Uint8(v[i])
		}
	case *[][]byte:
		slice = make([]uint8, len(*v))
		for i := range *v {
			slice[i] = Uint8((*v)[i])
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
			slice = make([]uint8, count)
			for i := 0; i < count; i++ {
				slice[i] = Uint8(rv.Index(i).Interface())
			}
		}
	}

	return
}

func Uint8Pointer(data any) *uint8 {
	v := Uint8(data)
	return &v
}

func Uint8sPointer(data any) *[]uint8 {
	v := Uint8s(data)
	return &v
}
