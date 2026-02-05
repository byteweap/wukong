package pulse

import "errors"

var (
	ErrClosed       = errors.New("pulse is closed")
	ErrBackpressure = errors.New("pulse backpressure: send queue full")
)
