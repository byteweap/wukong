package conv

func Uint16(data any) uint16 {
	return uint16(Uint64(data))
}

func Uint16s(data any) []uint16 {
	return convertSlice(data, Uint16)
}

func Uint16Pointer(data any) *uint16 {
	v := Uint16(data)
	return &v
}

func Uint16sPointer(data any) *[]uint16 {
	v := Uint16s(data)
	return &v
}
