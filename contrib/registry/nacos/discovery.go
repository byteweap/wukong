package nacos

import (
	"context"

	"github.com/byteweap/wukong/plugin/registry"
)

type Discovery struct {
}

var _ registry.Discovery = (*Discovery)(nil)

func (d Discovery) ID() string {
	//TODO implement me
	panic("implement me")
}

func (d Discovery) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	//TODO implement me
	panic("implement me")
}

func (d Discovery) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	//TODO implement me
	panic("implement me")
}
