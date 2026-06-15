package async

import (
	"sync"
	"testing"
	"time"
)

func TestGo_NormalExecution(t *testing.T) {
	done := make(chan struct{})
	Go(func() {
		close(done)
	}, nil)
	select {
	case <-done:
		// success
	case <-time.After(time.Second):
		t.Error("timeout waiting for function execution")
	}
}

func TestGo_NilFunction(_ *testing.T) {
	// Should not panic
	Go(nil, nil)
}

func TestGo_PanicWithHandler(t *testing.T) {
	recovered := make(chan any, 1)
	Go(func() {
		panic("test panic")
	}, func(r any) {
		recovered <- r
	})
	select {
	case r := <-recovered:
		if r != "test panic" {
			t.Errorf("expected recovered value 'test panic', got %v", r)
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for panic recovery")
	}
}

func TestGo_PanicWithNilHandler(t *testing.T) {
	// Should not call handler, but also should not crash
	done := make(chan struct{})
	Go(func() {
		defer func() {
			close(done)
		}()
		panic("test panic")
	}, nil)
	select {
	case <-done:
		// success
	case <-time.After(time.Second):
		t.Error("timeout waiting for goroutine completion")
	}
}

func TestRecover_WithPanicAndHandler(t *testing.T) {
	recovered := make(chan any, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer Recover(func(r any) {
			recovered <- r
		})
		panic("recover test")
	}()
	wg.Wait()
	select {
	case r := <-recovered:
		if r != "recover test" {
			t.Errorf("expected recovered value 'recover test', got %v", r)
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for panic recovery")
	}
}

func TestRecover_WithPanicAndNilHandler(_ *testing.T) {
	// Should not call handler, but also should not crash
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer Recover(nil)
		panic("recover test")
	}()
	wg.Wait()
	// No assertion needed, just ensure no crash
}

func TestRecover_NoPanic(t *testing.T) {
	var called bool
	Recover(func(_ any) {
		called = true
	})
	if called {
		t.Error("handler should not be called when no panic")
	}
}
