package conv

func Uint8(data any) uint8 {
	return uint8(Uint64(data))
}

func Uint8s(data any) []uint8 {
	return convertSlice(data, Uint8)
}

func Uint8Pointer(data any) *uint8 {
	v := Uint8(data)
	return &v
}

func Uint8sPointer(data any) *[]uint8 {
	v := Uint8s(data)
	return &v
}
