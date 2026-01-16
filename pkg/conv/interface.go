package conv

import "reflect"

func Interfaces(data any) (slice []any) {
	if data == nil {
		return
	}

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
		slice = make([]any, count)
		for i := 0; i < count; i++ {
			slice[i] = Int(rv.Index(i).Interface())
		}
	}

	return
}

func InterfacesPointer(data any) *[]any {
	v := Interfaces(data)
	return &v
}
