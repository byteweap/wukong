package conv

func Float32(data any) float32 {
	return float32(Float64(data))
}

func Float32s(data any) []float32 {
	return convertSlice(data, Float32)
}

func Float32Pointer(data any) *float32 {
	v := Float32(data)
	return &v
}

func Float32sPointer(data any) *[]float32 {
	v := Float32s(data)
	return &v
}
