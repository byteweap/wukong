package pulse

import (
	"time"

	"github.com/gobwas/ws"
)

type (
	OnClientOpenHandler    func(*ClientConn)
	OnClientCloseHandler   func(*ClientConn, error)
	OnClientMessageHandler func(*ClientConn, ws.OpCode, []byte)
	OnClientErrorHandler   func(*ClientConn, error)
)

type clientOptions struct {
	sendQueueSize        int
	maxMessageSize       int64
	readTimeout          time.Duration
	writeTimeout         time.Duration
	backpressure         BackpressureMode
	dialTimeout          time.Duration
	reconnectInterval    time.Duration
	reconnectMaxInterval time.Duration
	reconnectFactor      float64
	pingInterval         time.Duration
	pingTimeout          time.Duration

	onOpen    OnClientOpenHandler
	onClose   OnClientCloseHandler
	onMessage OnClientMessageHandler
	onError   OnClientErrorHandler
}

type ClientOption func(*clientOptions)

func defaultClientOptions() *clientOptions {
	return &clientOptions{
		sendQueueSize:        256,
		maxMessageSize:       64 * 1024,
		readTimeout:          0,
		writeTimeout:         0,
		backpressure:         BackpressureKick,
		dialTimeout:          10 * time.Second,
		reconnectInterval:    1 * time.Second,
		reconnectMaxInterval: 30 * time.Second,
		reconnectFactor:      2,
		pingInterval:         15 * time.Second,
		pingTimeout:          10 * time.Second,
	}
}

func ClientSendQueueSize(size int) ClientOption {
	return func(o *clientOptions) {
		o.sendQueueSize = size
	}
}

func ClientMaxMessageSize(size int64) ClientOption {
	return func(o *clientOptions) {
		o.maxMessageSize = size
	}
}

func ClientReadTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.readTimeout = timeout
	}
}

func ClientWriteTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.writeTimeout = timeout
	}
}

func ClientBackpressure(mode BackpressureMode) ClientOption {
	return func(o *clientOptions) {
		o.backpressure = mode
	}
}

func ClientDialTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.dialTimeout = timeout
	}
}

func ClientReconnectInterval(interval time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.reconnectInterval = interval
	}
}

func ClientReconnectMaxInterval(interval time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.reconnectMaxInterval = interval
	}
}

func ClientReconnectFactor(factor float64) ClientOption {
	return func(o *clientOptions) {
		o.reconnectFactor = factor
	}
}

func ClientPingInterval(interval time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.pingInterval = interval
	}
}

func ClientPingTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.pingTimeout = timeout
	}
}

func OnClientOpen(fn OnClientOpenHandler) ClientOption {
	return func(o *clientOptions) {
		o.onOpen = fn
	}
}

func OnClientClose(fn OnClientCloseHandler) ClientOption {
	return func(o *clientOptions) {
		o.onClose = fn
	}
}

func OnClientMessage(fn OnClientMessageHandler) ClientOption {
	return func(o *clientOptions) {
		o.onMessage = fn
	}
}

func OnClientError(fn OnClientErrorHandler) ClientOption {
	return func(o *clientOptions) {
		o.onError = fn
	}
}
