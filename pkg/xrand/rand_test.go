package xrand

import "testing"

func TestInt_EdgeCases(t *testing.T) {
	tests := []struct {
		name string
		min  int
		max  int
		want int
	}{
		{name: "max_zero", min: 5, max: 0, want: 0},
		{name: "min_eq_max", min: 7, max: 7, want: 7},
		{name: "min_gt_max", min: 9, max: 3, want: 3},
		{name: "neg_min_max_zero", min: -5, max: 0, want: 0},
		{name: "range_size_one_low", min: 0, max: 1, want: 0},
		{name: "range_size_one_high", min: 1, max: 2, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Int(tt.min, tt.max)
			if got != tt.want {
				t.Fatalf("Int(%d, %d) = %d, want %d", tt.min, tt.max, got, tt.want)
			}
		})
	}
}

func TestInt_Range(t *testing.T) {
	minValue, maxValue := -10, 10
	if minValue >= maxValue {
		t.Fatalf("invalid test range")
	}

	for i := 0; i < 1000; i++ {
		got := Int(minValue, maxValue)
		if got < minValue || got >= maxValue {
			t.Fatalf("Int(%d, %d) = %d out of range", minValue, maxValue, got)
		}
	}
}

func TestInt64_EdgeCases(t *testing.T) {
	tests := []struct {
		name string
		min  int64
		max  int64
		want int64
	}{
		{name: "max_zero", min: 5, max: 0, want: 0},
		{name: "min_eq_max", min: 7, max: 7, want: 7},
		{name: "min_gt_max", min: 9, max: 3, want: 3},
		{name: "neg_min_max_zero", min: -5, max: 0, want: 0},
		{name: "range_size_one_low", min: 0, max: 1, want: 0},
		{name: "range_size_one_high", min: 1, max: 2, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Int64(tt.min, tt.max)
			if got != tt.want {
				t.Fatalf("Int64(%d, %d) = %d, want %d", tt.min, tt.max, got, tt.want)
			}
		})
	}
}

func TestInt64_Range(t *testing.T) {
	var minValue int64 = -10
	var maxValue int64 = 10
	if minValue >= maxValue {
		t.Fatalf("invalid test range")
	}

	for i := 0; i < 1000; i++ {
		got := Int64(minValue, maxValue)
		if got < minValue || got >= maxValue {
			t.Fatalf("Int64(%d, %d) = %d out of range", minValue, maxValue, got)
		}
	}
}

func BenchmarkInt(b *testing.B) {
	minValue, maxValue := 0, 1_000_000
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Int(minValue, maxValue)
	}
}

func BenchmarkInt64(b *testing.B) {
	var minValue int64
	var maxValue int64 = 1_000_000
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Int64(minValue, maxValue)
	}
}
