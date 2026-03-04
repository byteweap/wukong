package server

import (
	"context"
	"net/url"
)

type Kind string

const (
	KindGate Kind = "gate"
	KindMesh Kind = "mesh"
)

type Server interface {
	Kind() Kind
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Endpoint(ctx context.Context) (*url.URL, error)
}
