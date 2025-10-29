package wos

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// WaitSignal 等待系统信号
func WaitSignal() {
	sig := make(chan os.Signal)
	switch runtime.GOOS {
	case `windows`:
		signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	default:
		signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM)
	}
	<-sig
	signal.Stop(sig)
}
