package conv

func Int32(data any) int32 {
	return int32(Int64(data))
}

func Int32s(data any) []int32 {
	return convertSlice(data, Int32)
}

func Int32Pointer(data any) *int32 {
	v := Int32(data)
	return &v
}

func Int32sPointer(data any) *[]int32 {
	v := Int32s(data)
	return &v
}
