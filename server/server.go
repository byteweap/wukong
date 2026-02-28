package server

import (
	"context"
	"net/url"

	"github.com/byteweap/wukong/internal/cluster"
)

type Server interface {
	Kind() cluster.Kind
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Endpoint() (*url.URL, error)
}
