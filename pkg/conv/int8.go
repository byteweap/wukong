package conv

func Int8(data any) int8 {
	return int8(Int64(data))
}

func Int8s(data any) []int8 {
	return convertSlice(data, Int8)
}

func Int8Pointer(data any) *int8 {
	v := Int8(data)
	return &v
}

func Int8sPointer(data any) *[]int8 {
	v := Int8s(data)
	return &v
}
