package selector

import "errors"

var (
	// ErrServiceUnavailable 表示服务不可用
	ErrServiceUnavailable = errors.New("service_unavailable")
	// ErrGatewayTimeout 表示网关超时
	ErrGatewayTimeout = errors.New("gateway_timeout")
)

// IsServiceUnavailable 判断是否服务不可用
func IsServiceUnavailable(err error) bool {
	return errors.Is(err, ErrServiceUnavailable)
}

// IsGatewayTimeout 判断是否网关超时
func IsGatewayTimeout(err error) bool {
	return errors.Is(err, ErrGatewayTimeout)
}
