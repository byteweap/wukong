package registry

import (
	"testing"
)

func TestServiceInstance_String(t *testing.T) {
	instance := &ServiceInstance{
		ID:   "123",
		Name: "test-service",
	}
	expected := "test-service-123"
	if instance.String() != expected {
		t.Errorf("expected %s, got %s", expected, instance.String())
	}
}

func TestServiceInstance_Equal(t *testing.T) {
	t.Run("both nil", func(t *testing.T) {
		var i *ServiceInstance
		if !i.Equal(nil) {
			t.Error("expected true for both nil")
		}
	})

	t.Run("first nil", func(t *testing.T) {
		var i *ServiceInstance
		o := &ServiceInstance{ID: "1"}
		if i.Equal(o) {
			t.Error("expected false for first nil")
		}
	})

	t.Run("second nil", func(t *testing.T) {
		i := &ServiceInstance{ID: "1"}
		if i.Equal(nil) {
			t.Error("expected false for second nil")
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		i := &ServiceInstance{ID: "1"}
		if i.Equal("not a ServiceInstance") {
			t.Error("expected false for wrong type")
		}
	})

	t.Run("different endpoints length", func(t *testing.T) {
		i := &ServiceInstance{
			ID:        "1",
			Endpoints: []string{"a"},
		}
		o := &ServiceInstance{
			ID:        "1",
			Endpoints: []string{"a", "b"},
		}
		if i.Equal(o) {
			t.Error("expected false for different endpoints length")
		}
	})

	t.Run("different endpoints", func(t *testing.T) {
		i := &ServiceInstance{
			ID:        "1",
			Endpoints: []string{"a"},
		}
		o := &ServiceInstance{
			ID:        "1",
			Endpoints: []string{"b"},
		}
		if i.Equal(o) {
			t.Error("expected false for different endpoints")
		}
	})

	t.Run("different metadata length", func(t *testing.T) {
		i := &ServiceInstance{
			ID:       "1",
			Metadata: map[string]string{"a": "1"},
		}
		o := &ServiceInstance{
			ID:       "1",
			Metadata: map[string]string{"a": "1", "b": "2"},
		}
		if i.Equal(o) {
			t.Error("expected false for different metadata length")
		}
	})

	t.Run("different metadata values", func(t *testing.T) {
		i := &ServiceInstance{
			ID:       "1",
			Metadata: map[string]string{"a": "1"},
		}
		o := &ServiceInstance{
			ID:       "1",
			Metadata: map[string]string{"a": "2"},
		}
		if i.Equal(o) {
			t.Error("expected false for different metadata values")
		}
	})

	t.Run("equal instances", func(t *testing.T) {
		i := &ServiceInstance{
			ID:        "1",
			Name:      "service",
			Version:   "v1",
			Endpoints: []string{"b", "a"},
			Metadata:  map[string]string{"key": "value"},
		}
		o := &ServiceInstance{
			ID:        "1",
			Name:      "service",
			Version:   "v1",
			Endpoints: []string{"a", "b"},
			Metadata:  map[string]string{"key": "value"},
		}
		if !i.Equal(o) {
			t.Error("expected true for equal instances")
		}
	})
}
