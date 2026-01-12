package nacos

import (
	"context"

	"github.com/byteweap/wukong/component/registry"
)

type Registry struct {
}

var _ registry.Registrar = (*Registry)(nil)

func (r *Registry) ID() string {
	//TODO implement me
	panic("implement me")
}

func (r *Registry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	//TODO implement me
	panic("implement me")
}

func (r *Registry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	//TODO implement me
	panic("implement me")
}
