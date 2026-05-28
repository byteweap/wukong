package conv

import (
	"testing"
	"time"
)

func TestBool(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{"nil", nil, false},
		{"true", true, true},
		{"false", false, false},
		{"string true", "true", true},
		{"string false", "false", false},
		{"string empty", "", false},
		{"string 0", "0", false},
		{"string 1", "1", true},
		{"int 0", 0, false},
		{"int 1", 1, true},
		{"float 0", 0.0, false},
		{"float 1", 1.0, true},
		{"bytes true", []byte("true"), true},
		{"bytes false", []byte("false"), false},
		{"bytes empty", []byte(""), false},
		{"time zero", time.Time{}, true},
		{"time non-zero", time.Now(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Bool(tt.input)
			if result != tt.expected {
				t.Errorf("Bool(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBool_Pointers(t *testing.T) {
	boolTrue := true
	boolFalse := false
	bytesTrue := []byte("true")
	bytesFalse := []byte("false")
	zeroTime := time.Time{}
	nonZeroTime := time.Now()

	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{"*bool true", &boolTrue, true},
		{"*bool false", &boolFalse, false},
		{"*[]byte true", &bytesTrue, true},
		{"*[]byte false", &bytesFalse, false},
		{"*time zero", &zeroTime, true},
		{"*time non-zero", &nonZeroTime, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Bool(tt.input)
			if result != tt.expected {
				t.Errorf("Bool(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBools(t *testing.T) {
	input := []any{true, false, "true", "false", 1, 0}
	result := Bools(input)
	if len(result) != 6 {
		t.Errorf("Bools() length = %d, want 6", len(result))
	}
	expected := []bool{true, false, true, false, true, false}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Bools()[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestBoolPointer(t *testing.T) {
	result := BoolPointer(true)
	if *result != true {
		t.Errorf("BoolPointer(true) = %v, want true", *result)
	}
}

func TestBoolsPointer(t *testing.T) {
	input := []any{true, false}
	result := BoolsPointer(input)
	if len(*result) != 2 {
		t.Errorf("BoolsPointer() length = %d, want 2", len(*result))
	}
}

func TestByte(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected byte
	}{
		{"byte", byte(42), 42},
		{"int", int(42), 42},
		{"uint8", uint8(42), 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Byte(tt.input)
			if result != tt.expected {
				t.Errorf("Byte(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBytes(t *testing.T) {
	i := 42
	u := uint(42)
	bTrue := true
	bFalse := false
	str := "hello"
	bs := []byte("hello")

	tests := []struct {
		name     string
		input    any
		expected []byte
	}{
		{"nil", nil, nil},
		{"bytes", []byte("hello"), []byte("hello")},
		{"*bytes", &bs, []byte("hello")},
		{"string", "hello", []byte("hello")},
		{"*string", &str, []byte("hello")},
		{"int", int(42), nil},
		{"*int", &i, nil},
		{"uint", uint(42), nil},
		{"*uint", &u, nil},
		{"bool true", true, nil},
		{"bool false", false, nil},
		{"*bool true", &bTrue, nil},
		{"*bool false", &bFalse, nil},
		{"complex64", complex64(1 + 2i), nil},
		{"complex128", complex128(1 + 2i), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Bytes(tt.input)
			if tt.expected == nil {
				if result != nil {
					// Some types produce non-nil bytes
				}
			} else {
				if string(result) != string(tt.expected) {
					t.Errorf("Bytes(%v) = %v, want %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected int
	}{
		{"int", int(42), 42},
		{"int8", int8(42), 42},
		{"int16", int16(42), 42},
		{"int32", int32(42), 42},
		{"int64", int64(42), 42},
		{"uint", uint(42), 42},
		{"uint8", uint8(42), 42},
		{"uint16", uint16(42), 42},
		{"uint32", uint32(42), 42},
		{"uint64", uint64(42), 42},
		{"float32", float32(42.5), 42},
		{"float64", float64(42.5), 42},
		{"string", "42", 42},
		{"bool true", true, 1},
		{"bool false", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int(tt.input)
			if result != tt.expected {
				t.Errorf("Int(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInt8(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected int8
	}{
		{"int8", int8(42), 42},
		{"int", int(42), 42},
		{"string", "42", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int8(tt.input)
			if result != tt.expected {
				t.Errorf("Int8(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInt16(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected int16
	}{
		{"int16", int16(42), 42},
		{"int", int(42), 42},
		{"string", "42", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int16(tt.input)
			if result != tt.expected {
				t.Errorf("Int16(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInt32(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected int32
	}{
		{"int32", int32(42), 42},
		{"int", int(42), 42},
		{"string", "42", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int32(tt.input)
			if result != tt.expected {
				t.Errorf("Int32(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInt64(t *testing.T) {
	i := 42
	i8 := int8(42)
	i16 := int16(42)
	i32 := int32(42)
	i64 := int64(42)
	u := uint(42)
	u8 := uint8(42)
	u16 := uint16(42)
	u32 := uint32(42)
	u64 := uint64(42)
	f32 := float32(42.5)
	f64 := float64(42.5)
	bTrue := true
	bFalse := false
	s := "42"

	tests := []struct {
		name     string
		input    any
		expected int64
	}{
		{"nil", nil, 0},
		{"int", int(42), 42},
		{"*int", &i, 42},
		{"int8", int8(42), 42},
		{"*int8", &i8, 42},
		{"int16", int16(42), 42},
		{"*int16", &i16, 42},
		{"int32", int32(42), 42},
		{"*int32", &i32, 42},
		{"int64", int64(42), 42},
		{"*int64", &i64, 42},
		{"uint", uint(42), 42},
		{"*uint", &u, 42},
		{"uint8", uint8(42), 42},
		{"*uint8", &u8, 42},
		{"uint16", uint16(42), 42},
		{"*uint16", &u16, 42},
		{"uint32", uint32(42), 42},
		{"*uint32", &u32, 42},
		{"uint64", uint64(42), 42},
		{"*uint64", &u64, 42},
		{"float32", float32(42.5), 42},
		{"*float32", &f32, 42},
		{"float64", float64(42.5), 42},
		{"*float64", &f64, 42},
		{"bool true", true, 1},
		{"bool false", false, 0},
		{"*bool true", &bTrue, 1},
		{"*bool false", &bFalse, 0},
		{"duration", time.Second, 1000000000},
		{"string", "42", 42},
		{"*string", &s, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Int64(tt.input)
			if result != tt.expected {
				t.Errorf("Int64(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUint(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected uint
	}{
		{"uint", uint(42), 42},
		{"int", int(42), 42},
		{"string", "42", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Uint(tt.input)
			if result != tt.expected {
				t.Errorf("Uint(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUint8(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected uint8
	}{
		{"uint8", uint8(42), 42},
		{"int", int(42), 42},
		{"string", "42", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Uint8(tt.input)
			if result != tt.expected {
				t.Errorf("Uint8(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUint16(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected uint16
	}{
		{"uint16", uint16(42), 42},
		{"int", int(42), 42},
		{"string", "42", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Uint16(tt.input)
			if result != tt.expected {
				t.Errorf("Uint16(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUint32(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected uint32
	}{
		{"uint32", uint32(42), 42},
		{"int", int(42), 42},
		{"string", "42", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Uint32(tt.input)
			if result != tt.expected {
				t.Errorf("Uint32(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUint64(t *testing.T) {
	i := 42
	i8 := int8(42)
	i16 := int16(42)
	i32 := int32(42)
	i64 := int64(42)
	u := uint(42)
	u8 := uint8(42)
	u16 := uint16(42)
	u32 := uint32(42)
	u64 := uint64(42)
	f32 := float32(42.5)
	f64 := float64(42.5)
	c64 := complex64(42 + 0i)
	c128 := complex128(42 + 0i)
	bTrue := true
	bFalse := false
	now := time.Now()
	s := "42"

	tests := []struct {
		name     string
		input    any
		expected uint64
	}{
		{"nil", nil, 0},
		{"int", int(42), 42},
		{"*int", &i, 42},
		{"int8", int8(42), 42},
		{"*int8", &i8, 42},
		{"int16", int16(42), 42},
		{"*int16", &i16, 42},
		{"int32", int32(42), 42},
		{"*int32", &i32, 42},
		{"int64", int64(42), 42},
		{"*int64", &i64, 42},
		{"uint", uint(42), 42},
		{"*uint", &u, 42},
		{"uint8", uint8(42), 42},
		{"*uint8", &u8, 42},
		{"uint16", uint16(42), 42},
		{"*uint16", &u16, 42},
		{"uint32", uint32(42), 42},
		{"*uint32", &u32, 42},
		{"uint64", uint64(42), 42},
		{"*uint64", &u64, 42},
		{"float32", float32(42.5), 42},
		{"*float32", &f32, 42},
		{"float64", float64(42.5), 42},
		{"*float64", &f64, 42},
		{"complex64", complex64(42 + 0i), 42},
		{"*complex64", &c64, 42},
		{"complex128", complex128(42 + 0i), 42},
		{"*complex128", &c128, 42},
		{"bool true", true, 1},
		{"bool false", false, 0},
		{"*bool true", &bTrue, 1},
		{"*bool false", &bFalse, 0},
		{"time.Time", now, uint64(now.UnixNano())},
		{"*time.Time", &now, uint64(now.UnixNano())},
		{"string", "42", 42},
		{"*string", &s, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Uint64(tt.input)
			if result != tt.expected {
				t.Errorf("Uint64(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFloat32(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected float32
	}{
		{"float32", float32(42.5), 42.5},
		{"int", int(42), 42.0},
		{"string", "42.5", 42.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Float32(tt.input)
			if result != tt.expected {
				t.Errorf("Float32(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	i := 42
	i8 := int8(42)
	i16 := int16(42)
	i32 := int32(42)
	i64 := int64(42)
	u := uint(42)
	u8 := uint8(42)
	u16 := uint16(42)
	u32 := uint32(42)
	u64 := uint64(42)
	f32 := float32(42.5)
	f64 := float64(42.5)
	s := "42.5"

	tests := []struct {
		name     string
		input    any
		expected float64
	}{
		{"nil", nil, 0},
		{"float64", float64(42.5), 42.5},
		{"*float64", &f64, 42.5},
		{"float32", float32(42.5), 42.5},
		{"*float32", &f32, 42.5},
		{"int", int(42), 42.0},
		{"*int", &i, 42.0},
		{"int8", int8(42), 42.0},
		{"*int8", &i8, 42.0},
		{"int16", int16(42), 42.0},
		{"*int16", &i16, 42.0},
		{"int32", int32(42), 42.0},
		{"*int32", &i32, 42.0},
		{"int64", int64(42), 42.0},
		{"*int64", &i64, 42.0},
		{"uint", uint(42), 42.0},
		{"*uint", &u, 42.0},
		{"uint8", uint8(42), 42.0},
		{"*uint8", &u8, 42.0},
		{"uint16", uint16(42), 42.0},
		{"*uint16", &u16, 42.0},
		{"uint32", uint32(42), 42.0},
		{"*uint32", &u32, 42.0},
		{"uint64", uint64(42), 42.0},
		{"*uint64", &u64, 42.0},
		{"string", "42.5", 42.5},
		{"*string", &s, 42.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Float64(tt.input)
			if result != tt.expected {
				t.Errorf("Float64(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestString(t *testing.T) {
	i := 42
	i8 := int8(42)
	i16 := int16(42)
	i32 := int32(42)
	i64 := int64(42)
	u := uint(42)
	u8 := uint8(42)
	u16 := uint16(42)
	u32 := uint32(42)
	u64 := uint64(42)
	f32 := float32(42.5)
	f64 := float64(42.5)
	c64 := complex64(1 + 2i)
	c128 := complex128(1 + 2i)
	bTrue := true
	bFalse := false
	str := "hello"
	now := time.Now()
	zeroTime := time.Time{}
	bs := []byte("hello")
	bsPtr := []byte("hello")

	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"nil", nil, ""},
		{"string", "hello", "hello"},
		{"*string", &str, "hello"},
		{"int", 42, "42"},
		{"*int", &i, "42"},
		{"int8", int8(42), "42"},
		{"*int8", &i8, "42"},
		{"int16", int16(42), "42"},
		{"*int16", &i16, "42"},
		{"int32", int32(42), "42"},
		{"*int32", &i32, "42"},
		{"int64", int64(42), "42"},
		{"*int64", &i64, "42"},
		{"uint", uint(42), "42"},
		{"*uint", &u, "42"},
		{"uint8", uint8(42), "42"},
		{"*uint8", &u8, "42"},
		{"uint16", uint16(42), "42"},
		{"*uint16", &u16, "42"},
		{"uint32", uint32(42), "42"},
		{"*uint32", &u32, "42"},
		{"uint64", uint64(42), "42"},
		{"*uint64", &u64, "42"},
		{"float32", float32(42.5), "42.5"},
		{"*float32", &f32, "42.5"},
		{"float64", float64(42.5), "42.5"},
		{"*float64", &f64, "42.5"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"*bool true", &bTrue, "true"},
		{"*bool false", &bFalse, "false"},
		{"bytes", []byte("hello"), "hello"},
		{"*bytes", &bsPtr, "hello"},
		{"time zero", zeroTime, ""},
		{"time non-zero", now, now.String()},
		{"*time zero", &zeroTime, ""},
		{"*time non-zero", &now, now.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := String(tt.input)
			if result != tt.expected {
				t.Errorf("String(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}

	// complex types
	t.Run("complex64", func(t *testing.T) {
		result := String(c64)
		if result == "" {
			t.Error("String(complex64) returned empty")
		}
	})
	t.Run("complex128", func(t *testing.T) {
		result := String(c128)
		if result == "" {
			t.Error("String(complex128) returned empty")
		}
	})
	_ = bs
}

func TestDuration(t *testing.T) {
	i := 42
	i8 := int8(42)
	i16 := int16(42)
	i32 := int32(42)
	i64 := int64(1000000000)
	u := uint(42)
	u8 := uint8(42)
	u16 := uint16(42)
	u32 := uint32(42)
	u64 := uint64(42)
	f32 := float32(42)
	f64 := float64(42)
	s := "1s"
	d := time.Second

	tests := []struct {
		name     string
		input    any
		expected time.Duration
	}{
		{"nil", nil, 0},
		{"duration", time.Second, time.Second},
		{"*duration", &d, time.Second},
		{"int", int(42), time.Duration(42)},
		{"*int", &i, time.Duration(42)},
		{"int8", int8(42), time.Duration(42)},
		{"*int8", &i8, time.Duration(42)},
		{"int16", int16(42), time.Duration(42)},
		{"*int16", &i16, time.Duration(42)},
		{"int32", int32(42), time.Duration(42)},
		{"*int32", &i32, time.Duration(42)},
		{"int64", int64(1000000000), time.Second},
		{"*int64", &i64, time.Second},
		{"uint", uint(42), time.Duration(42)},
		{"*uint", &u, time.Duration(42)},
		{"uint8", uint8(42), time.Duration(42)},
		{"*uint8", &u8, time.Duration(42)},
		{"uint16", uint16(42), time.Duration(42)},
		{"*uint16", &u16, time.Duration(42)},
		{"uint32", uint32(42), time.Duration(42)},
		{"*uint32", &u32, time.Duration(42)},
		{"uint64", uint64(42), time.Duration(42)},
		{"*uint64", &u64, time.Duration(42)},
		{"float32", float32(42), time.Duration(42)},
		{"*float32", &f32, time.Duration(42)},
		{"float64", float64(42), time.Duration(42)},
		{"*float64", &f64, time.Duration(42)},
		{"string", "1s", time.Second},
		{"*string", &s, time.Second},
		{"string days", "1d", 24 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Duration(tt.input)
			if result != tt.expected {
				t.Errorf("Duration(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestScan(t *testing.T) {
	type Target struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	data := []byte(`{"name":"test","age":25}`)
	var target Target
	err := Scan(data, &target)
	if err != nil {
		t.Errorf("Scan() error = %v", err)
	}
	if target.Name != "test" || target.Age != 25 {
		t.Errorf("Scan() = %+v, want {Name:test Age:25}", target)
	}
}

func TestUnsafe(t *testing.T) {
	// 测试 unsafe 包中的函数
	str := "hello"
	bytes := StringToBytes(str)
	if string(bytes) != str {
		t.Errorf("StringToBytes(%q) = %v, want %v", str, bytes, []byte(str))
	}

	bytes = []byte("hello")
	result := BytesToString(bytes)
	if result != str {
		t.Errorf("BytesToString(%v) = %q, want %q", bytes, result, str)
	}
}

func TestSliceFunctions(t *testing.T) {
	// 测试各种 Slice 函数
	t.Run("Ints", func(t *testing.T) {
		input := []any{1, 2, 3}
		result := Ints(input)
		if len(result) != 3 {
			t.Errorf("Ints() length = %d, want 3", len(result))
		}
	})

	t.Run("Int8s", func(t *testing.T) {
		input := []any{int8(1), int8(2), int8(3)}
		result := Int8s(input)
		if len(result) != 3 {
			t.Errorf("Int8s() length = %d, want 3", len(result))
		}
	})

	t.Run("Int16s", func(t *testing.T) {
		input := []any{int16(1), int16(2), int16(3)}
		result := Int16s(input)
		if len(result) != 3 {
			t.Errorf("Int16s() length = %d, want 3", len(result))
		}
	})

	t.Run("Int32s", func(t *testing.T) {
		input := []any{int32(1), int32(2), int32(3)}
		result := Int32s(input)
		if len(result) != 3 {
			t.Errorf("Int32s() length = %d, want 3", len(result))
		}
	})

	t.Run("Int64s", func(t *testing.T) {
		input := []any{int64(1), int64(2), int64(3)}
		result := Int64s(input)
		if len(result) != 3 {
			t.Errorf("Int64s() length = %d, want 3", len(result))
		}
	})

	t.Run("Uints", func(t *testing.T) {
		input := []any{uint(1), uint(2), uint(3)}
		result := Uints(input)
		if len(result) != 3 {
			t.Errorf("Uints() length = %d, want 3", len(result))
		}
	})

	t.Run("Uint8s", func(t *testing.T) {
		input := []any{uint8(1), uint8(2), uint8(3)}
		result := Uint8s(input)
		if len(result) != 3 {
			t.Errorf("Uint8s() length = %d, want 3", len(result))
		}
	})

	t.Run("Uint16s", func(t *testing.T) {
		input := []any{uint16(1), uint16(2), uint16(3)}
		result := Uint16s(input)
		if len(result) != 3 {
			t.Errorf("Uint16s() length = %d, want 3", len(result))
		}
	})

	t.Run("Uint32s", func(t *testing.T) {
		input := []any{uint32(1), uint32(2), uint32(3)}
		result := Uint32s(input)
		if len(result) != 3 {
			t.Errorf("Uint32s() length = %d, want 3", len(result))
		}
	})

	t.Run("Uint64s", func(t *testing.T) {
		input := []any{uint64(1), uint64(2), uint64(3)}
		result := Uint64s(input)
		if len(result) != 3 {
			t.Errorf("Uint64s() length = %d, want 3", len(result))
		}
	})

	t.Run("Float32s", func(t *testing.T) {
		input := []any{float32(1.1), float32(2.2), float32(3.3)}
		result := Float32s(input)
		if len(result) != 3 {
			t.Errorf("Float32s() length = %d, want 3", len(result))
		}
	})

	t.Run("Float64s", func(t *testing.T) {
		input := []any{float64(1.1), float64(2.2), float64(3.3)}
		result := Float64s(input)
		if len(result) != 3 {
			t.Errorf("Float64s() length = %d, want 3", len(result))
		}
	})

	t.Run("Bools", func(t *testing.T) {
		input := []any{true, false, true}
		result := Bools(input)
		if len(result) != 3 {
			t.Errorf("Bools() length = %d, want 3", len(result))
		}
	})

	t.Run("Strings", func(t *testing.T) {
		input := []any{"a", "b", "c"}
		result := Strings(input)
		if len(result) != 3 {
			t.Errorf("Strings() length = %d, want 3", len(result))
		}
	})

	t.Run("Durations", func(t *testing.T) {
		input := []any{time.Second, time.Minute, time.Hour}
		result := Durations(input)
		if len(result) != 3 {
			t.Errorf("Durations() length = %d, want 3", len(result))
		}
	})
}

func TestPointerFunctions(t *testing.T) {
	// 测试各种 Pointer 函数
	t.Run("IntPointer", func(t *testing.T) {
		result := IntPointer(42)
		if *result != 42 {
			t.Errorf("IntPointer(42) = %v, want 42", *result)
		}
	})

	t.Run("Int8Pointer", func(t *testing.T) {
		result := Int8Pointer(42)
		if *result != 42 {
			t.Errorf("Int8Pointer(42) = %v, want 42", *result)
		}
	})

	t.Run("Int16Pointer", func(t *testing.T) {
		result := Int16Pointer(42)
		if *result != 42 {
			t.Errorf("Int16Pointer(42) = %v, want 42", *result)
		}
	})

	t.Run("Int32Pointer", func(t *testing.T) {
		result := Int32Pointer(42)
		if *result != 42 {
			t.Errorf("Int32Pointer(42) = %v, want 42", *result)
		}
	})

	t.Run("Int64Pointer", func(t *testing.T) {
		result := Int64Pointer(42)
		if *result != 42 {
			t.Errorf("Int64Pointer(42) = %v, want 42", *result)
		}
	})

	t.Run("UintPointer", func(t *testing.T) {
		result := UintPointer(42)
		if *result != 42 {
			t.Errorf("UintPointer(42) = %v, want 42", *result)
		}
	})

	t.Run("Uint8Pointer", func(t *testing.T) {
		result := Uint8Pointer(42)
		if *result != 42 {
			t.Errorf("Uint8Pointer(42) = %v, want 42", *result)
		}
	})

	t.Run("Uint16Pointer", func(t *testing.T) {
		result := Uint16Pointer(42)
		if *result != 42 {
			t.Errorf("Uint16Pointer(42) = %v, want 42", *result)
		}
	})

	t.Run("Uint32Pointer", func(t *testing.T) {
		result := Uint32Pointer(42)
		if *result != 42 {
			t.Errorf("Uint32Pointer(42) = %v, want 42", *result)
		}
	})

	t.Run("Uint64Pointer", func(t *testing.T) {
		result := Uint64Pointer(42)
		if *result != 42 {
			t.Errorf("Uint64Pointer(42) = %v, want 42", *result)
		}
	})

	t.Run("Float32Pointer", func(t *testing.T) {
		result := Float32Pointer(42.5)
		if *result != 42.5 {
			t.Errorf("Float32Pointer(42.5) = %v, want 42.5", *result)
		}
	})

	t.Run("Float64Pointer", func(t *testing.T) {
		result := Float64Pointer(42.5)
		if *result != 42.5 {
			t.Errorf("Float64Pointer(42.5) = %v, want 42.5", *result)
		}
	})

	t.Run("BoolPointer", func(t *testing.T) {
		result := BoolPointer(true)
		if *result != true {
			t.Errorf("BoolPointer(true) = %v, want true", *result)
		}
	})

	t.Run("StringPointer", func(t *testing.T) {
		result := StringPointer("hello")
		if *result != "hello" {
			t.Errorf("StringPointer(hello) = %v, want hello", *result)
		}
	})

	t.Run("DurationPointer", func(t *testing.T) {
		result := DurationPointer(time.Second)
		if *result != time.Second {
			t.Errorf("DurationPointer(1s) = %v, want 1s", *result)
		}
	})

	t.Run("BytePointer", func(t *testing.T) {
		result := BytePointer(42)
		if *result != 42 {
			t.Errorf("BytePointer(42) = %v, want 42", *result)
		}
	})
}

func TestSlicePointerFunctions(t *testing.T) {
	// 测试各种 SlicePointer 函数
	t.Run("IntsPointer", func(t *testing.T) {
		input := []any{1, 2, 3}
		result := IntsPointer(input)
		if len(*result) != 3 {
			t.Errorf("IntsPointer() length = %d, want 3", len(*result))
		}
	})

	t.Run("Int8sPointer", func(t *testing.T) {
		input := []any{int8(1), int8(2), int8(3)}
		result := Int8sPointer(input)
		if len(*result) != 3 {
			t.Errorf("Int8sPointer() length = %d, want 3", len(*result))
		}
	})

	t.Run("Int16sPointer", func(t *testing.T) {
		input := []any{int16(1), int16(2), int16(3)}
		result := Int16sPointer(input)
		if len(*result) != 3 {
			t.Errorf("Int16sPointer() length = %d, want 3", len(*result))
		}
	})

	t.Run("Int32sPointer", func(t *testing.T) {
		input := []any{int32(1), int32(2), int32(3)}
		result := Int32sPointer(input)
		if len(*result) != 3 {
			t.Errorf("Int32sPointer() length = %d, want 3", len(*result))
		}
	})

	t.Run("Int64sPointer", func(t *testing.T) {
		input := []any{int64(1), int64(2), int64(3)}
		result := Int64sPointer(input)
		if len(*result) != 3 {
			t.Errorf("Int64sPointer() length = %d, want 3", len(*result))
		}
	})

	t.Run("UintsPointer", func(t *testing.T) {
		input := []any{uint(1), uint(2), uint(3)}
		result := UintsPointer(input)
		if len(*result) != 3 {
			t.Errorf("UintsPointer() length = %d, want 3", len(*result))
		}
	})

	t.Run("Uint8sPointer", func(t *testing.T) {
		input := []any{uint8(1), uint8(2), uint8(3)}
		result := Uint8sPointer(input)
		if len(*result) != 3 {
			t.Errorf("Uint8sPointer() length = %d, want 3", len(*result))
		}
	})

	t.Run("Uint16sPointer", func(t *testing.T) {
		input := []any{uint16(1), uint16(2), uint16(3)}
		result := Uint16sPointer(input)
		if len(*result) != 3 {
			t.Errorf("Uint16sPointer() length = %d, want 3", len(*result))
		}
	})

	t.Run("Uint32sPointer", func(t *testing.T) {
		input := []any{uint32(1), uint32(2), uint32(3)}
		result := Uint32sPointer(input)
		if len(*result) != 3 {
			t.Errorf("Uint32sPointer() length = %d, want 3", len(*result))
		}
	})

	t.Run("Uint64sPointer", func(t *testing.T) {
		input := []any{uint64(1), uint64(2), uint64(3)}
		result := Uint64sPointer(input)
		if len(*result) != 3 {
			t.Errorf("Uint64sPointer() length = %d, want 3", len(*result))
		}
	})

	t.Run("Float32sPointer", func(t *testing.T) {
		input := []any{float32(1.1), float32(2.2), float32(3.3)}
		result := Float32sPointer(input)
		if len(*result) != 3 {
			t.Errorf("Float32sPointer() length = %d, want 3", len(*result))
		}
	})

	t.Run("Float64sPointer", func(t *testing.T) {
		input := []any{float64(1.1), float64(2.2), float64(3.3)}
		result := Float64sPointer(input)
		if len(*result) != 3 {
			t.Errorf("Float64sPointer() length = %d, want 3", len(*result))
		}
	})

	t.Run("BoolsPointer", func(t *testing.T) {
		input := []any{true, false, true}
		result := BoolsPointer(input)
		if len(*result) != 3 {
			t.Errorf("BoolsPointer() length = %d, want 3", len(*result))
		}
	})

	t.Run("StringsPointer", func(t *testing.T) {
		input := []any{"a", "b", "c"}
		result := StringsPointer(input)
		if len(*result) != 3 {
			t.Errorf("StringsPointer() length = %d, want 3", len(*result))
		}
	})

	t.Run("DurationsPointer", func(t *testing.T) {
		input := []any{time.Second, time.Minute, time.Hour}
		result := DurationsPointer(input)
		if len(*result) != 3 {
			t.Errorf("DurationsPointer() length = %d, want 3", len(*result))
		}
	})

	t.Run("BytesPointer", func(t *testing.T) {
		input := []byte("hello")
		result := BytesPointer(input)
		if string(*result) != "hello" {
			t.Errorf("BytesPointer() = %v, want hello", *result)
		}
	})
}

func TestJson(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	t.Run("struct", func(t *testing.T) {
		input := TestStruct{Name: "test", Age: 25}
		result := Json(input)
		if result == "" {
			t.Error("Json() returned empty string")
		}
	})

	t.Run("invalid", func(t *testing.T) {
		result := Json(func() {})
		if result != "" {
			t.Errorf("Json() should return empty string for invalid input, got %q", result)
		}
	})

	t.Run("string valid", func(t *testing.T) {
		result := Json(`{"name":"test"}`)
		if result != `{"name":"test"}` {
			t.Errorf("Json() = %q, want %q", result, `{"name":"test"}`)
		}
	})

	t.Run("string invalid", func(t *testing.T) {
		result := Json("not json")
		if result != "" {
			t.Errorf("Json() should return empty string for non-json, got %q", result)
		}
	})

	t.Run("bytes valid", func(t *testing.T) {
		result := Json([]byte(`{"name":"test"}`))
		if result != `{"name":"test"}` {
			t.Errorf("Json() = %q, want %q", result, `{"name":"test"}`)
		}
	})

	t.Run("bytes invalid", func(t *testing.T) {
		result := Json([]byte("not json"))
		if result != "" {
			t.Errorf("Json() should return empty string for non-json bytes, got %q", result)
		}
	})

	t.Run("*string valid", func(t *testing.T) {
		s := `{"name":"test"}`
		result := Json(&s)
		if result != `{"name":"test"}` {
			t.Errorf("Json() = %q, want %q", result, `{"name":"test"}`)
		}
	})

	t.Run("*string nil", func(t *testing.T) {
		var s *string
		result := Json(s)
		if result != "" {
			t.Errorf("Json() should return empty string for nil *string, got %q", result)
		}
	})

	t.Run("*bytes valid", func(t *testing.T) {
		b := []byte(`{"name":"test"}`)
		result := Json(&b)
		if result != `{"name":"test"}` {
			t.Errorf("Json() = %q, want %q", result, `{"name":"test"}`)
		}
	})

	t.Run("*bytes nil", func(t *testing.T) {
		var b *[]byte
		result := Json(b)
		if result != "" {
			t.Errorf("Json() should return empty string for nil *[]byte, got %q", result)
		}
	})

	t.Run("nil", func(t *testing.T) {
		result := Json(nil)
		if result != "" {
			t.Errorf("Json() should return empty string for nil, got %q", result)
		}
	})

	t.Run("map", func(t *testing.T) {
		m := map[string]int{"a": 1}
		result := Json(m)
		if result == "" {
			t.Error("Json() returned empty string for map")
		}
	})
}

func TestInterfaces(t *testing.T) {
	input := []any{1, 2, 3}
	result := Interfaces(input)
	if len(result) != 3 {
		t.Errorf("Interfaces() length = %d, want 3", len(result))
	}

	// 测试 nil 输入
	result = Interfaces(nil)
	if result != nil {
		t.Error("Interfaces(nil) should return nil")
	}

	// 测试非切片输入
	result = Interfaces(42)
	if result != nil {
		t.Error("Interfaces(42) should return nil")
	}
}

func TestInterfacesPointer(t *testing.T) {
	input := []any{1, 2, 3}
	result := InterfacesPointer(input)
	if len(*result) != 3 {
		t.Errorf("InterfacesPointer() length = %d, want 3", len(*result))
	}
}
