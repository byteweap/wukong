package kconv

import "reflect"

func Int32(data any) int32 {
	return int32(Int64(data))
}

func Int32s(data any) (slice []int32) {

	if data == nil {
		return
	}
	switch v := data.(type) {
	case []int:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]int:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case []int8:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]int8:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case []int16:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]int16:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case []int32:
		return v
	case *[]int32:
		return *v
	case []int64:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]int64:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case []uint:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]uint:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case []uint8:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]uint8:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case []uint16:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]uint16:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case []uint32:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]uint32:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case []uint64:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]uint64:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case []float32:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]float32:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case []float64:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]float64:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case []complex64:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]complex64:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case []complex128:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]complex128:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case []string:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]string:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case []bool:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]bool:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case []any:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[]any:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
		}
	case [][]byte:
		slice = make([]int32, len(v))
		for i := range v {
			slice[i] = Int32(v[i])
		}
	case *[][]byte:
		slice = make([]int32, len(*v))
		for i := range *v {
			slice[i] = Int32((*v)[i])
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
			slice = make([]int32, count)
			for i := 0; i < count; i++ {
				slice[i] = Int32(rv.Index(i).Interface())
			}
		}
	}

	return
}

func Int32Pointer(data any) *int32 {
	v := Int32(data)
	return &v
}

func Int32sPointer(data any) *[]int32 {
	v := Int32s(data)
	return &v
}
