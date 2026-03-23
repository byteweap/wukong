package errors

import "errors"

var (
	ErrOptionsRequired   = errors.New("options required")
	ErrAppNotFound       = errors.New("app not found")
	ErrLocatorRequired   = errors.New("locator required")
	ErrBrokerRequired    = errors.New("broker required")
	ErrDiscoveryRequired = errors.New("discovery required")
	ErrSelectorRequired  = errors.New("selector func required")
)
