package kconv

import "reflect"

func Float32(data any) float32 {
	return float32(Float64(data))
}

func Float32s(data any) (slice []float32) {
	if data == nil {
		return
	}

	switch v := data.(type) {
	case []int:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]int:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case []int8:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]int8:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case []int16:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]int16:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case []int32:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]int32:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case []int64:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]int64:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case []uint:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]uint:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case []uint8:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]uint8:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case []uint16:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]uint16:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case []uint32:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]uint32:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case []uint64:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]uint64:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case []float32:
		return v
	case *[]float32:
		return *v
	case []float64:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]float64:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case []complex64:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]complex64:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case []complex128:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]complex128:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case []string:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]string:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case []bool:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]bool:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case []any:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[]any:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
		}
	case [][]byte:
		slice = make([]float32, len(v))
		for i := range v {
			slice[i] = Float32(v[i])
		}
	case *[][]byte:
		slice = make([]float32, len(*v))
		for i := range *v {
			slice[i] = Float32((*v)[i])
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
			slice = make([]float32, count)
			for i := 0; i < count; i++ {
				slice[i] = Float32(rv.Index(i).Interface())
			}
		}
	}

	return
}

func Float32Pointer(data any) *float32 {
	v := Float32(data)
	return &v
}

func Float32sPointer(data any) *[]float32 {
	v := Float32s(data)
	return &v
}
