package conv

func Uint(data any) uint {
	return uint(Uint64(data))
}

// Uints 任何类型转uint切片
func Uints(data any) []uint {
	return convertSlice(data, Uint)
}

func UintPointer(data any) *uint {
	v := Uint(data)
	return &v
}

func UintsPointer(data any) *[]uint {
	v := Uints(data)
	return &v
}
