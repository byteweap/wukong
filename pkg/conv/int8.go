package conv

import "reflect"

func Int8(data any) int8 {
	return int8(Int64(data))
}

func Int8s(data any) (slice []int8) {
	if data == nil {
		return
	}

	switch v := data.(type) {
	case []int:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]int:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case []int8:
		return v
	case *[]int8:
		return *v
	case []int16:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]int16:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case []int32:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]int32:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case []int64:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]int64:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case []uint:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]uint:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case []uint8:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]uint8:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case []uint16:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]uint16:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case []uint32:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]uint32:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case []uint64:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]uint64:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case []float32:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]float32:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case []float64:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]float64:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case []complex64:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]complex64:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case []complex128:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]complex128:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case []string:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]string:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case []bool:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]bool:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case []any:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[]any:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
		}
	case [][]byte:
		slice = make([]int8, len(v))
		for i := range v {
			slice[i] = Int8(v[i])
		}
	case *[][]byte:
		slice = make([]int8, len(*v))
		for i := range *v {
			slice[i] = Int8((*v)[i])
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
			slice = make([]int8, count)
			for i := 0; i < count; i++ {
				slice[i] = Int8(rv.Index(i).Interface())
			}
		}
	}

	return
}

func Int8Pointer(data any) *int8 {
	v := Int8(data)
	return &v
}

func Int8sPointer(data any) *[]int8 {
	v := Int8s(data)
	return &v
}
