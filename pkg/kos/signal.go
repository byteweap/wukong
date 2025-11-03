package kos

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// WaitSignal 等待系统信号
func WaitSignal() {
	sign := make(chan os.Signal)
	switch runtime.GOOS {
	case `windows`:
		signal.Notify(sign, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	default:
		signal.Notify(sign, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM)
	}
	<-sign
	signal.Stop(sign)
}
