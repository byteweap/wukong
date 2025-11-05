package kconv

import "reflect"

func Int(data any) int {
	return int(Int64(data))
}

func Ints(data any) (slice []int) {
	if data == nil {
		return
	}

	switch v := data.(type) {
	case []int:
		return v
	case *[]int:
		return *v
	case []int8:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]int8:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case []int16:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]int16:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case []int32:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]int32:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case []int64:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]int64:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case []uint:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]uint:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case []uint8:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]uint8:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case []uint16:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]uint16:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case []uint32:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]uint32:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case []uint64:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]uint64:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case []float32:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]float32:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case []float64:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]float64:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case []complex64:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]complex64:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case []complex128:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]complex128:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case []string:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]string:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case []bool:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]bool:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case []any:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[]any:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
		}
	case [][]byte:
		slice = make([]int, len(v))
		for i := range v {
			slice[i] = Int(v[i])
		}
	case *[][]byte:
		slice = make([]int, len(*v))
		for i := range *v {
			slice[i] = Int((*v)[i])
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
			slice = make([]int, count)
			for i := 0; i < count; i++ {
				slice[i] = Int(rv.Index(i).Interface())
			}
		}
	}

	return
}

func IntPointer(data any) *int {
	v := Int(data)
	return &v
}

func IntsPointer(data any) *[]int {
	v := Ints(data)
	return &v
}
