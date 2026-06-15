package conv

import (
	"encoding/json"
	"reflect"
)

func JSON(data any) string {
	switch v := data.(type) {
	case string:
		if isRawJSON(v) {
			return v
		}
		return ""
	case *string:
		if v != nil && isRawJSON(*v) {
			return *v
		}
		return ""
	case []byte:
		if s := BytesToString(v); isRawJSON(s) {
			return s
		}
		return ""
	case *[]byte:
		if v != nil {
			s := BytesToString(*v)
			if isRawJSON(s) {
				return s
			}
		}
		return ""
	}

	rv := indirectValue(data)
	if !rv.IsValid() {
		return ""
	}

	switch rv.Kind() {
	case reflect.String:
		if s := rv.String(); isRawJSON(s) {
			return s
		}
	case reflect.Map, reflect.Array, reflect.Slice, reflect.Struct:
		if b, err := json.Marshal(data); err == nil {
			return BytesToString(b)
		}
	default:
		return ""
	}

	return ""
}

func isRawJSON(s string) bool {
	l := len(s)
	return l >= 2 && ((s[0] == '{' && s[l-1] == '}') || (s[0] == '[' && s[l-1] == ']'))
}
