package lang

import (
	"testing"
)

func TestIf(t *testing.T) {
	// Test with bool
	t.Run("bool true", func(t *testing.T) {
		result := If(true, true, false)
		if result != true {
			t.Errorf("If(true, true, false) = %v, want true", result)
		}
	})

	t.Run("bool false", func(t *testing.T) {
		result := If(false, true, false)
		if result != false {
			t.Errorf("If(false, true, false) = %v, want false", result)
		}
	})

	// Test with int
	t.Run("int true", func(t *testing.T) {
		result := If(true, 1, 2)
		if result != 1 {
			t.Errorf("If(true, 1, 2) = %v, want 1", result)
		}
	})

	t.Run("int false", func(t *testing.T) {
		result := If(false, 1, 2)
		if result != 2 {
			t.Errorf("If(false, 1, 2) = %v, want 2", result)
		}
	})

	// Test with string
	t.Run("string true", func(t *testing.T) {
		result := If(true, "yes", "no")
		if result != "yes" {
			t.Errorf("If(true, 'yes', 'no') = %v, want 'yes'", result)
		}
	})

	t.Run("string false", func(t *testing.T) {
		result := If(false, "yes", "no")
		if result != "no" {
			t.Errorf("If(false, 'yes', 'no') = %v, want 'no'", result)
		}
	})

	// Test with float64
	t.Run("float64 true", func(t *testing.T) {
		result := If(true, 1.5, 2.5)
		if result != 1.5 {
			t.Errorf("If(true, 1.5, 2.5) = %v, want 1.5", result)
		}
	})

	t.Run("float64 false", func(t *testing.T) {
		result := If(false, 1.5, 2.5)
		if result != 2.5 {
			t.Errorf("If(false, 1.5, 2.5) = %v, want 2.5", result)
		}
	})
}
