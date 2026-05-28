package conv

func Int(data any) int {
	return int(Int64(data))
}

func Ints(data any) []int {
	return convertSlice(data, Int)
}

func IntPointer(data any) *int {
	v := Int(data)
	return &v
}

func IntsPointer(data any) *[]int {
	v := Ints(data)
	return &v
}
