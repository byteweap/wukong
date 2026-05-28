package conv

import "reflect"

func indirectValue(data any) reflect.Value {
	rv := reflect.ValueOf(data)
	for rv.IsValid() && rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return rv
}

func convertSlice[T any](data any, convert func(any) T) []T {
	if data == nil {
		return nil
	}

	switch v := data.(type) {
	case []T:
		return v
	case *[]T:
		return *v
	}

	rv := indirectValue(data)

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		count := rv.Len()
		slice := make([]T, count)
		for i := 0; i < count; i++ {
			slice[i] = convert(rv.Index(i).Interface())
		}
		return slice
	default:
		return nil
	}
}
