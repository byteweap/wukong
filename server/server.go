package server

import (
	"context"
	"net/url"
)

type Server interface {
	Kind() Kind
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Endpoint() (*url.URL, error)
}

type Kind string

const (
	KindGate Kind = "gate"
	KindGame Kind = "game"
)
