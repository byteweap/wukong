package wos

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// getSignalsForOS returns signals to monitor based on the operating system
func getSignalsForOS() []os.Signal {
	switch runtime.GOOS {
	case "windows":
		// Windows supported signals
		return []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	default:
		// Unix/Linux supported signals (excluding SIGKILL as it cannot be caught)
		return []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT}
	}
}

// WaitSignal blocks until a system signal is received
func WaitSignal() {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, getSignalsForOS()...)
	<-sign
	signal.Stop(sign)
}

// Signal creates and returns a signal channel
// Callers can decide how to handle the signals
func Signal() <-chan os.Signal {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, getSignalsForOS()...)
	return sign
}
