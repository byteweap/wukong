package polaris

import (
	"context"

	"github.com/byteweap/wukong/component/registry"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/model"
)

type watcher struct {
	ServiceName      string
	Namespace        string
	Ctx              context.Context
	Cancel           context.CancelFunc
	Channel          <-chan model.SubScribeEvent
	ServiceInstances []*registry.ServiceInstance
	first            bool
}

var _ registry.Watcher = (*watcher)(nil)

func newWatcher(ctx context.Context, namespace string, serviceName string, consumer api.ConsumerAPI) (*watcher, error) {
	watchServiceResponse, err := consumer.WatchService(&api.WatchServiceRequest{
		WatchServiceRequest: model.WatchServiceRequest{
			Key: model.ServiceKey{
				Namespace: namespace,
				Service:   serviceName,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	w := &watcher{
		Namespace:        namespace,
		ServiceName:      serviceName,
		first:            true,
		Channel:          watchServiceResponse.EventChannel,
		ServiceInstances: instancesToServiceInstances(watchServiceResponse.GetAllInstancesResp.GetInstances()),
	}
	w.Ctx, w.Cancel = context.WithCancel(ctx)
	return w, nil
}

// Next 在以下情况返回实例列表
// 1 首次 watch 且实例列表不为空
// 2 实例发生变更
// 否则阻塞直到 ctx 超时或取消
func (w *watcher) Next() ([]*registry.ServiceInstance, error) {
	if w.first {
		w.first = false
		return w.ServiceInstances, nil
	}
	select {
	case <-w.Ctx.Done():
		return nil, w.Ctx.Err()
	case event := <-w.Channel:
		if event.GetSubScribeEventType() == model.EventInstance {
			// 这里通常为真，但仍需检查，防止事件类型变化
			if instanceEvent, ok := event.(*model.InstanceEvent); ok {
				// 处理 DeleteEvent
				if instanceEvent.DeleteEvent != nil {
					for _, instance := range instanceEvent.DeleteEvent.Instances {
						for i, serviceInstance := range w.ServiceInstances {
							if serviceInstance.ID == instance.GetId() {
								// 移除匹配项
								if len(w.ServiceInstances) <= 1 {
									w.ServiceInstances = w.ServiceInstances[:0]
									continue
								}
								w.ServiceInstances = append(w.ServiceInstances[:i], w.ServiceInstances[i+1:]...)
							}
						}
					}
				}
				// 处理 UpdateEvent
				if instanceEvent.UpdateEvent != nil {
					for i, serviceInstance := range w.ServiceInstances {
						for _, update := range instanceEvent.UpdateEvent.UpdateList {
							if serviceInstance.ID == update.Before.GetId() {
								w.ServiceInstances[i] = instanceToServiceInstance(update.After)
							}
						}
					}
				}
				// 处理 AddEvent
				if instanceEvent.AddEvent != nil {
					w.ServiceInstances = append(w.ServiceInstances, instancesToServiceInstances(instanceEvent.AddEvent.Instances)...)
				}
			}
			return w.ServiceInstances, nil
		}
	}
	return w.ServiceInstances, nil
}

// Stop 关闭 watcher
func (w *watcher) Stop() error {
	w.Cancel()
	return nil
}
