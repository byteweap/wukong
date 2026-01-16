package conv

import "reflect"

func Uint(data any) uint {
	return uint(Uint64(data))
}

// Uints 任何类型转uint切片
func Uints(data any) (slice []uint) {
	if data == nil {
		return
	}

	switch v := data.(type) {
	case []int:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]int:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []int8:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]int8:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []int16:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]int16:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []int32:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]int32:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []int64:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]int64:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []uint:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]uint:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []uint8:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]uint8:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []uint16:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]uint16:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []uint32:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]uint32:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []uint64:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]uint64:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []float32:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]float32:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []float64:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]float64:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []complex64:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]complex64:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []complex128:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]complex128:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []string:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]string:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []bool:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]bool:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case []any:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[]any:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
		}
	case [][]byte:
		slice = make([]uint, len(v))
		for i := range v {
			slice[i] = Uint(v[i])
		}
	case *[][]byte:
		slice = make([]uint, len(*v))
		for i := range *v {
			slice[i] = Uint((*v)[i])
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
			slice = make([]uint, count)
			for i := 0; i < count; i++ {
				slice[i] = Uint(rv.Index(i).Interface())
			}
		}
	}

	return
}

func UintPointer(data any) *uint {
	v := Uint(data)
	return &v
}

func UintsPointer(data any) *[]uint {
	v := Uints(data)
	return &v
}
