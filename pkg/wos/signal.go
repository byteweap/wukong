package wos

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// getSignalsForOS 根据操作系统返回需要监听的信号
func getSignalsForOS() []os.Signal {
	switch runtime.GOOS {
	case "windows":
		// Windows 支持的信号
		return []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	default:
		// Unix/Linux 支持的信号（不包括 SIGKILL，因为它无法被捕获）
		return []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT}
	}
}

// WaitSignal 阻塞直到收到系统信号
func WaitSignal() {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, getSignalsForOS()...)
	<-sign
	signal.Stop(sign)
}

// Signal 创建并返回信号通道
// 调用者可以自行决定如何处理信号
func Signal() <-chan os.Signal {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, getSignalsForOS()...)
	return sign
}
