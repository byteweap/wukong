package etcd

import "github.com/byteweap/wukong/plugin/registry"

type Watcher struct {
}

var _ registry.Watcher = (*Watcher)(nil)

func (w *Watcher) ID() string {
	//TODO implement me
	panic("implement me")
}

func (w *Watcher) Next() ([]*registry.ServiceInstance, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Watcher) Stop() error {
	//TODO implement me
	panic("implement me")
}
