package conv

func Uint32(data any) uint32 {
	return uint32(Uint64(data))
}

func Uint32s(data any) []uint32 {
	return convertSlice(data, Uint32)
}

func Uint32Pointer(data any) *uint32 {
	v := Uint32(data)
	return &v
}

func Uint32sPointer(data any) *[]uint32 {
	v := Uint32s(data)
	return &v
}
