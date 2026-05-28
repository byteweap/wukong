package conv

func Int16(data any) int16 {
	return int16(Int64(data))
}

func Int16s(data any) []int16 {
	return convertSlice(data, Int16)
}

func Int16Pointer(data any) *int16 {
	v := Int16(data)
	return &v
}

func Int16sPointer(data any) *[]int16 {
	v := Int16s(data)
	return &v
}
