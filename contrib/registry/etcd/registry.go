package etcd

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/byteweap/wukong/component/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Registry 是 etcd 注册中心实现
type Registry struct {
	opts   *options
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
	// ctxMap 保存每个服务实例的取消函数
	// 注销服务时调用对应的取消函数以停止心跳
	ctxMap map[string]*serviceCancel
}

var _ registry.Registry = (*Registry)(nil)

func (r *Registry) ID() string {
	return "etcd"
}

type serviceCancel struct {
	service *registry.ServiceInstance
	cancel  context.CancelFunc
}

// New 创建 etcd 注册中心
func New(client *clientv3.Client, opts ...Option) (r *Registry) {
	op := &options{
		ctx:       context.Background(),
		namespace: "/microservices",
		ttl:       time.Second * 15,
		maxRetry:  5,
	}
	for _, o := range opts {
		o(op)
	}
	return &Registry{
		opts:   op,
		client: client,
		kv:     clientv3.NewKV(client),
		ctxMap: make(map[string]*serviceCancel),
	}
}

// Register 注册服务
func (r *Registry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	key := r.registerKey(service)
	value, err := marshal(service)
	if err != nil {
		return err
	}
	if r.lease != nil {
		r.lease.Close()
	}
	r.lease = clientv3.NewLease(r.client)
	leaseID, err := r.registerWithKV(ctx, key, value)
	if err != nil {
		return err
	}

	hctx, cancel := context.WithCancel(r.opts.ctx)
	r.ctxMap[key] = &serviceCancel{
		service: service,
		cancel:  cancel,
	}
	go r.heartBeat(hctx, leaseID, key, value)
	return nil
}

func (r *Registry) registerKey(service *registry.ServiceInstance) string {
	return fmt.Sprintf("%s/%s/%s", r.opts.namespace, service.Name, service.ID)
}

// Deregister 注销服务
func (r *Registry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	defer func() {
		if r.lease != nil {
			r.lease.Close()
		}
	}()
	// 取消心跳
	key := r.registerKey(service)
	if serviceCancel, ok := r.ctxMap[key]; ok {
		serviceCancel.cancel()
		delete(r.ctxMap, key)
	}
	_, err := r.client.Delete(ctx, key)
	return err
}

// GetService 按服务名获取实例列表
func (r *Registry) GetService(ctx context.Context, name string) ([]*registry.ServiceInstance, error) {
	key := r.serviceKey(name)
	resp, err := r.kv.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	items := make([]*registry.ServiceInstance, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		si, err := unmarshal(kv.Value)
		if err != nil {
			return nil, err
		}
		if si.Name != name {
			continue
		}
		items = append(items, si)
	}
	return items, nil
}

func (r *Registry) serviceKey(name string) string {
	return fmt.Sprintf("%s/%s", r.opts.namespace, name)
}

// Watch 按服务名创建 watcher
func (r *Registry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	key := r.serviceKey(name)
	return newWatcher(ctx, key, name, r.client)
}

// registerWithKV 创建租约并写入 KV，返回租约 ID
func (r *Registry) registerWithKV(ctx context.Context, key string, value string) (clientv3.LeaseID, error) {
	grant, err := r.lease.Grant(ctx, int64(r.opts.ttl.Seconds()))
	if err != nil {
		return 0, err
	}
	_, err = r.client.Put(ctx, key, value, clientv3.WithLease(grant.ID))
	if err != nil {
		return 0, err
	}
	return grant.ID, nil
}

func (r *Registry) heartBeat(ctx context.Context, leaseID clientv3.LeaseID, key string, value string) {
	curLeaseID := leaseID
	kac, err := r.client.KeepAlive(ctx, leaseID)
	if err != nil {
		curLeaseID = 0
	}
	randSource := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))

	for {
		if curLeaseID == 0 {
			// 重新尝试注册
			var retreat []int
			for retryCnt := 0; retryCnt < r.opts.maxRetry; retryCnt++ {
				if ctx.Err() != nil {
					return
				}
				// 防止无限阻塞
				idChan := make(chan clientv3.LeaseID, 1)
				errChan := make(chan error, 1)
				cancelCtx, cancel := context.WithCancel(ctx)
				go func() {
					defer cancel()
					id, registerErr := r.registerWithKV(cancelCtx, key, value)
					if registerErr != nil {
						errChan <- registerErr
					} else {
						idChan <- id
					}
				}()

				select {
				case <-time.After(3 * time.Second):
					cancel()
					continue
				case <-errChan:
					continue
				case curLeaseID = <-idChan:
				}

				kac, err = r.client.KeepAlive(ctx, curLeaseID)
				if err == nil {
					break
				}
				retreat = append(retreat, 1<<retryCnt)
				time.Sleep(time.Duration(retreat[randSource.IntN(len(retreat))]) * time.Second)
			}
			if _, ok := <-kac; !ok {
				// 重试失败
				return
			}
		}

		select {
		case _, ok := <-kac:
			if !ok {
				if ctx.Err() != nil {
					// 上下文取消导致通道关闭
					return
				}
				// 需要重新注册
				curLeaseID = 0
				continue
			}
		case <-r.opts.ctx.Done():
			return
		}
	}
}
