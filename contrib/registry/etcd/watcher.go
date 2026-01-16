package etcd

import (
	"context"
	"fmt"
	"sync"

	"github.com/byteweap/wukong/component/registry"
	"github.com/byteweap/wukong/encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Watcher 使用 etcd Watch 实现服务监听
type Watcher struct {
	client    *clientv3.Client
	prefix    string
	namespace string
	ctx       context.Context
	cancel    context.CancelFunc
	watchCh   clientv3.WatchChan
	eventCh   chan []*registry.ServiceInstance
	errCh     chan error
	mu        sync.RWMutex
	instances map[string]*registry.ServiceInstance // key -> instance
	once      sync.Once
	stopped   bool
}

var _ registry.Watcher = (*Watcher)(nil)

// newWatcher 创建新的监听器
func newWatcher(ctx context.Context, client *clientv3.Client, prefix, namespace string) (*Watcher, error) {

	watchCtx, cancel := context.WithCancel(ctx)

	w := &Watcher{
		client:    client,
		prefix:    prefix,
		namespace: namespace,
		ctx:       watchCtx,
		cancel:    cancel,
		watchCh:   client.Watch(watchCtx, prefix, clientv3.WithPrefix()),
		eventCh:   make(chan []*registry.ServiceInstance, 1),
		errCh:     make(chan error, 1),
		instances: make(map[string]*registry.ServiceInstance),
	}

	// 首次获取当前所有实例
	if err := w.initialLoad(ctx); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initial load services: %w", err)
	}

	// 启动 watch goroutine
	go w.watchLoop()

	return w, nil
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

		// 取消 context，停止 watch
		w.cancel()

		// 关闭 channel
		close(w.eventCh)
		close(w.errCh)
	})
	return err
}

// initialLoad 首次加载所有服务实例
func (w *Watcher) initialLoad(ctx context.Context) error {
	resp, err := w.client.Get(ctx, w.prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	for _, kv := range resp.Kvs {
		var instance registry.ServiceInstance
		if err := json.Unmarshal(kv.Value, &instance); err != nil {
			continue
		}
		w.instances[string(kv.Key)] = &instance
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

// watchLoop 监听 etcd 变更事件
func (w *Watcher) watchLoop() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case resp, ok := <-w.watchCh:
			if !ok {
				// channel 关闭
				w.mu.RLock()
				stopped := w.stopped
				w.mu.RUnlock()
				if !stopped {
					select {
					case w.errCh <- fmt.Errorf("watch channel closed"):
					default:
					}
				}
				return
			}

			if resp.Err() != nil {
				select {
				case w.errCh <- resp.Err():
				default:
				}
				continue
			}

			// 处理事件
			if w.processEvents(resp.Events) {
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
		}
	}
}

// processEvents 处理 watch 事件，返回是否有变更
func (w *Watcher) processEvents(events []*clientv3.Event) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	changed := false
	for _, event := range events {
		key := string(event.Kv.Key)

		switch event.Type {
		case clientv3.EventTypePut:
			// 新增或更新
			var instance registry.ServiceInstance
			if err := json.Unmarshal(event.Kv.Value, &instance); err != nil {
				continue
			}
			// 检查是否有变更
			old, exists := w.instances[key]
			if !exists || !old.Equal(&instance) {
				w.instances[key] = &instance
				changed = true
			}

		case clientv3.EventTypeDelete:
			// 删除
			if _, exists := w.instances[key]; exists {
				delete(w.instances, key)
				changed = true
			}
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
