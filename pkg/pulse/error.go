package pulse

import "errors"

var (
	ErrClosed         = errors.New("pulse server is closed")
	ErrBackpressure   = errors.New("pulse backpressure: send queue full")
	ErrNotConnected   = errors.New("pulse client not connected")
	ErrAlreadyRunning = errors.New("pulse client already running")
	ErrPingTimeout    = errors.New("pulse client ping timeout")
)
