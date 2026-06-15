package wos

import (
	"os"
	"runtime"
	"syscall"
	"testing"
	"time"
)

func TestGetSignalsForOS(t *testing.T) {
	signals := getSignalsForOS()
	if len(signals) == 0 {
		t.Error("getSignalsForOS() returned empty slice")
	}

	switch runtime.GOOS {
	case windowsOS:
		if len(signals) != 2 {
			t.Errorf("getSignalsForOS() for windows returned %d signals, want 2", len(signals))
		}
	default:
		if len(signals) != 4 {
			t.Errorf("getSignalsForOS() for unix returned %d signals, want 4", len(signals))
		}
	}
}

func TestSignal(t *testing.T) {
	ch := Signal()
	if ch == nil {
		t.Error("Signal() returned nil channel")
	}

	// Test that channel is of correct type
	assertSignalChannel(ch)
}

func assertSignalChannel(_ <-chan os.Signal) {}

func TestWaitSignal(t *testing.T) {
	if runtime.GOOS == windowsOS {
		t.Skip("os.Process.Signal does not support SIGINT on Windows")
	}

	done := make(chan struct{})
	go func() {
		WaitSignal()
		close(done)
	}()

	// Send a signal to ourselves
	time.Sleep(100 * time.Millisecond)
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Signal(syscall.SIGINT); err != nil {
		t.Fatal(err)
	}

	select {
	case <-done:
		// success
	case <-time.After(time.Second):
		t.Error("WaitSignal did not return after receiving signal")
	}
}
