package nacos

import (
	"context"
	"fmt"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/byteweap/wukong/component/registry"
)

// ID Nacos 监听器实现标识符
const WatcherID = "nacos(watcher)"

// Watcher 使用 Nacos Subscribe 实现服务监听
type Watcher struct {
	namingClient naming_client.INamingClient
	serviceName  string
	opts         *options
	ctx          context.Context
	cancel       context.CancelFunc
	eventCh      chan []*registry.ServiceInstance
	errCh        chan error
	mu           sync.RWMutex
	instances    map[string]*registry.ServiceInstance // instanceId -> instance
	once         sync.Once
	stopped      bool
}

var _ registry.Watcher = (*Watcher)(nil)

// newWatcher 创建新的监听器
func newWatcher(ctx context.Context, namingClient naming_client.INamingClient, serviceName string, opts *options) (*Watcher, error) {
	watchCtx, cancel := context.WithCancel(ctx)

	w := &Watcher{
		namingClient: namingClient,
		serviceName:  serviceName,
		opts:         opts,
		ctx:          watchCtx,
		cancel:       cancel,
		eventCh:      make(chan []*registry.ServiceInstance, 1),
		errCh:        make(chan error, 1),
		instances:    make(map[string]*registry.ServiceInstance),
	}

	// 首次加载当前所有实例
	if err := w.initialLoad(ctx); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initial load services: %w", err)
	}

	// 订阅服务变更
	if err := w.subscribe(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to subscribe service: %w", err)
	}

	return w, nil
}

// ID 返回实现标识符
func (w *Watcher) ID() string {
	return WatcherID
}

// Next 阻塞等待服务实例变更，首次调用或变更时返回实例列表
func (w *Watcher) Next() ([]*registry.ServiceInstance, error) {
	select {
	case instances := <-w.eventCh:
		if instances == nil {
			// channel 关闭
			return nil, fmt.Errorf("watcher is stopped")
		}
		return instances, nil
	case err := <-w.errCh:
		return nil, err
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	}
}

// Stop 停止监听并释放资源
func (w *Watcher) Stop() error {
	var err error
	w.once.Do(func() {
		w.mu.Lock()
		w.stopped = true
		w.mu.Unlock()

		// 取消 context，停止订阅
		w.cancel()

		// 取消订阅
		param := vo.SubscribeParam{
			ServiceName: w.serviceName,
			GroupName:   w.opts.group,
			Clusters:    []string{w.opts.clusterName},
		}
		_ = w.namingClient.Unsubscribe(&param)

		// 关闭 channel
		close(w.eventCh)
		close(w.errCh)
	})
	return err
}

// initialLoad 首次加载所有服务实例
func (w *Watcher) initialLoad(ctx context.Context) error {
	param := vo.SelectInstancesParam{
		ServiceName: w.serviceName,
		GroupName:   w.opts.group,
		Clusters:    []string{w.opts.clusterName},
		HealthyOnly: false,
	}

	instances, err := w.namingClient.SelectInstances(param)
	if err != nil {
		return err
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	for _, instance := range instances {
		si := convertToServiceInstance(instance, w.serviceName)
		if si != nil {
			w.instances[si.ID] = si
		}
	}

	// 如果有实例，立即发送
	if len(w.instances) > 0 {
		instances := w.getInstancesLocked()
		select {
		case w.eventCh <- instances:
		default:
			// channel 已满，跳过
		}
	}

	return nil
}

// subscribe 订阅服务变更
func (w *Watcher) subscribe() error {
	param := vo.SubscribeParam{
		ServiceName: w.serviceName,
		GroupName:   w.opts.group,
		Clusters:    []string{w.opts.clusterName},
		SubscribeCallback: func(services []model.Instance, err error) {
			if err != nil {
				w.mu.RLock()
				stopped := w.stopped
				w.mu.RUnlock()
				if !stopped {
					select {
					case w.errCh <- err:
					default:
					}
				}
				return
			}

			// 处理服务变更
			if w.processInstances(services) {
				// 有变更，发送最新实例列表
				w.mu.RLock()
				instances := w.getInstancesLocked()
				stopped := w.stopped
				w.mu.RUnlock()

				if !stopped {
					select {
					case w.eventCh <- instances:
					case <-w.ctx.Done():
						return
					}
				}
			}
		},
	}

	return w.namingClient.Subscribe(&param)
}

// processInstances 处理服务实例列表，返回是否有变更
func (w *Watcher) processInstances(services []model.Instance) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 构建新的实例映射
	newInstances := make(map[string]*registry.ServiceInstance)
	for _, instance := range services {
		si := convertToServiceInstance(instance, w.serviceName)
		if si != nil {
			newInstances[si.ID] = si
		}
	}

	// 检查是否有变更
	changed := false

	// 检查新增或更新的实例
	for id, newInstance := range newInstances {
		oldInstance, exists := w.instances[id]
		if !exists || !oldInstance.Equal(newInstance) {
			w.instances[id] = newInstance
			changed = true
		}
	}

	// 检查删除的实例
	for id := range w.instances {
		if _, exists := newInstances[id]; !exists {
			delete(w.instances, id)
			changed = true
		}
	}

	return changed
}

// getInstancesLocked 获取所有实例列表（需要已持有锁）
func (w *Watcher) getInstancesLocked() []*registry.ServiceInstance {
	instances := make([]*registry.ServiceInstance, 0, len(w.instances))
	for _, instance := range w.instances {
		instances = append(instances, instance)
	}
	return instances
}
