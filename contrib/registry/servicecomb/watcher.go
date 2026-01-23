package servicecomb

import (
	"context"

	"github.com/byteweap/wukong/component/registry"
	"github.com/go-chassis/sc-client"
)

type Watcher struct {
	cli RegistryClient
	ch  chan *registry.ServiceInstance
}

var _ registry.Watcher = (*Watcher)(nil)

func newWatcher(_ context.Context, cli RegistryClient, serviceName string) (*Watcher, error) {
	// 建立当前服务与目标服务的依赖关系，用于发现
	_, err := cli.FindMicroServiceInstances(curServiceID, appID, serviceName, "")
	if err != nil {
		return nil, err
	}
	w := &Watcher{
		cli: cli,
		ch:  make(chan *registry.ServiceInstance),
	}
	go func() {
		watchErr := w.cli.WatchMicroService(curServiceID, func(event *sc.MicroServiceInstanceChangedEvent) {
			if event.Key.ServiceName != serviceName {
				return
			}
			svcIns := &registry.ServiceInstance{
				ID:        event.Instance.InstanceId,
				Name:      event.Key.ServiceName,
				Version:   event.Key.Version,
				Metadata:  event.Instance.Properties,
				Endpoints: event.Instance.Endpoints,
			}
			w.Put(svcIns)
		})
		if watchErr != nil {
			return
		}
	}()
	return w, nil
}

// Put 仅用于单元测试
func (w *Watcher) Put(svcIns *registry.ServiceInstance) {
	w.ch <- svcIns
}

func (w *Watcher) Next() ([]*registry.ServiceInstance, error) {
	var svcInstances []*registry.ServiceInstance
	svcIns := <-w.ch
	svcInstances = append(svcInstances, svcIns)
	return svcInstances, nil
}

func (w *Watcher) Stop() error {
	close(w.ch)
	return nil
}
