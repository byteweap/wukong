package consul

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/byteweap/wukong/component/registry"
	"github.com/hashicorp/consul/api"
)

// Config 是 consul 注册中心配置
type Config struct {
	*api.Config
}

// Registry 是 consul 注册中心实现
type Registry struct {
	cli               *Client
	enableHealthCheck bool
	registry          map[string]*serviceSet
	lock              sync.RWMutex
	timeout           time.Duration
}

var _ registry.Registry = (*Registry)(nil)

// New 创建 consul 注册中心
func New(apiClient *api.Client, opts ...Option) *Registry {
	r := &Registry{
		registry:          make(map[string]*serviceSet),
		enableHealthCheck: true,
		timeout:           10 * time.Second,
		cli: &Client{
			dc:                             SingleDatacenter,
			cli:                            apiClient,
			resolver:                       defaultResolver,
			healthcheckInterval:            10,
			heartbeat:                      true,
			deregisterCriticalServiceAfter: 600,
			cancelers:                      make(map[string]*canceler),
		},
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

func (r *Registry) ID() string {
	return "consul"
}

// Register 注册服务
func (r *Registry) Register(ctx context.Context, svc *registry.ServiceInstance) error {
	return r.cli.Register(ctx, svc, r.enableHealthCheck)
}

// Deregister 注销服务
func (r *Registry) Deregister(ctx context.Context, svc *registry.ServiceInstance) error {
	return r.cli.Deregister(ctx, svc.ID)
}

// GetService 按名称获取服务
func (r *Registry) GetService(ctx context.Context, name string) ([]*registry.ServiceInstance, error) {
	r.lock.RLock()
	set := r.registry[name]
	r.lock.RUnlock()

	getRemote := func() []*registry.ServiceInstance {
		services, _, err := r.cli.Service(ctx, name, 0, true)
		if err == nil && len(services) > 0 {
			return services
		}
		return nil
	}

	if set == nil {
		if s := getRemote(); len(s) > 0 {
			return s, nil
		}
		return nil, fmt.Errorf("service %s not resolved in registry", name)
	}
	ss, _ := set.services.Load().([]*registry.ServiceInstance)
	if ss == nil {
		if s := getRemote(); len(s) > 0 {
			return s, nil
		}
		return nil, fmt.Errorf("service %s not found in registry", name)
	}
	return ss, nil
}

// ListServices 返回所有服务
func (r *Registry) ListServices() (allServices map[string][]*registry.ServiceInstance, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	allServices = make(map[string][]*registry.ServiceInstance)
	for name, set := range r.registry {
		var services []*registry.ServiceInstance
		ss, _ := set.services.Load().([]*registry.ServiceInstance)
		if ss == nil {
			continue
		}
		services = append(services, ss...)
		allServices[name] = services
	}
	return
}

// Watch 按名称监听服务变更
func (r *Registry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.lock.Lock()
	set, ok := r.registry[name]
	if !ok {
		cancelCtx, cancel := context.WithCancel(context.Background())
		set = &serviceSet{
			registry:    r,
			watcher:     make(map[*watcher]struct{}),
			services:    &atomic.Value{},
			serviceName: name,
			ctx:         cancelCtx,
			cancel:      cancel,
		}
		r.registry[name] = set
	}
	set.ref.Add(1)
	r.lock.Unlock()

	// 初始化 watcher
	w := &watcher{
		event: make(chan struct{}, 1),
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	w.set = set
	set.lock.Lock()
	set.watcher[w] = struct{}{}
	set.lock.Unlock()

	ss, _ := set.services.Load().([]*registry.ServiceInstance)
	if len(ss) > 0 {
		// 如果已有缓存数据，先推送一次，避免首次 watch 一直阻塞
		select {
		case w.event <- struct{}{}:
		default:
		}
	}

	if !ok {
		if err := r.resolve(ctx, set); err != nil {
			return nil, err
		}
	}
	return w, nil
}

// resolve 启动后台拉取并广播服务变更
func (r *Registry) resolve(ctx context.Context, ss *serviceSet) error {
	listServices := r.cli.Service
	if r.timeout > 0 {
		listServices = func(ctx context.Context, service string, index uint64, passingOnly bool) ([]*registry.ServiceInstance, uint64, error) {
			timeoutCtx, cancel := context.WithTimeout(ctx, r.timeout)
			defer cancel()

			return r.cli.Service(timeoutCtx, service, index, passingOnly)
		}
	}

	services, idx, err := listServices(ctx, ss.serviceName, 0, true)
	if err != nil {
		return err
	}
	if len(services) > 0 {
		ss.broadcast(services)
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				tmpService, tmpIdx, err := listServices(ss.ctx, ss.serviceName, idx, true)
				if err != nil {
					if err := sleepCtx(ss.ctx, time.Second); err != nil {
						return
					}
					continue
				}
				if len(tmpService) != 0 && tmpIdx != idx {
					services = tmpService
					ss.broadcast(services)
				}
				idx = tmpIdx
			case <-ss.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (r *Registry) tryDelete(ss *serviceSet) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	if ss.ref.Add(-1) != 0 {
		return false
	}
	ss.cancel()
	delete(r.registry, ss.serviceName)
	return true
}
