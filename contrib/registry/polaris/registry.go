package polaris

import (
	"context"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/component/registry"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"
	"github.com/polarismesh/polaris-go/pkg/model"
)

// _instanceIDSeparator 用于拼接实例 ID
const _instanceIDSeparator = "-"

// Registry 是 polaris 注册中心实现
type Registry struct {
	opt      options
	provider api.ProviderAPI
	consumer api.ConsumerAPI
}

var _ registry.Registry = (*Registry)(nil)

func NewRegistry(provider api.ProviderAPI, consumer api.ConsumerAPI, opts ...Option) (r *Registry) {
	op := options{
		Namespace:    "default",
		ServiceToken: "",
		Protocol:     nil,
		Weight:       0,
		Priority:     0,
		Healthy:      true,
		Heartbeat:    true,
		Isolate:      false,
		TTL:          0,
		Timeout:      0,
		RetryCount:   0,
	}
	for _, option := range opts {
		option(&op)
	}
	return &Registry{
		opt:      op,
		provider: provider,
		consumer: consumer,
	}
}

func NewRegistryWithConfig(conf config.Configuration, opts ...Option) (r *Registry) {
	provider, err := api.NewProviderAPIByConfig(conf)
	if err != nil {
		panic(err)
	}
	consumer, err := api.NewConsumerAPIByConfig(conf)
	if err != nil {
		panic(err)
	}
	return NewRegistry(provider, consumer, opts...)
}

func (r *Registry) ID() string {
	return "polaris"
}

// Register 注册服务
func (r *Registry) Register(_ context.Context, serviceInstance *registry.ServiceInstance) error {
	ids := make([]string, 0, len(serviceInstance.Endpoints))
	for _, endpoint := range serviceInstance.Endpoints {
		// 解析 url
		u, err := url.Parse(endpoint)
		if err != nil {
			return err
		}

		// 解析 host 与 port
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return err
		}

		// 端口转为整数
		portNum, err := strconv.Atoi(port)
		if err != nil {
			return err
		}

		// 组装 metadata
		var rmd map[string]string
		if serviceInstance.Metadata == nil {
			rmd = map[string]string{
				"kind":    u.Scheme,
				"version": serviceInstance.Version,
			}
		} else {
			rmd = make(map[string]string, len(serviceInstance.Metadata)+2)
			for k, v := range serviceInstance.Metadata {
				rmd[k] = v
			}
			rmd["kind"] = u.Scheme
			rmd["version"] = serviceInstance.Version
		}
		// 调用注册
		service, err := r.provider.Register(
			&api.InstanceRegisterRequest{
				InstanceRegisterRequest: model.InstanceRegisterRequest{
					Service:      serviceInstance.Name + u.Scheme,
					ServiceToken: r.opt.ServiceToken,
					Namespace:    r.opt.Namespace,
					Host:         host,
					Port:         portNum,
					Protocol:     r.opt.Protocol,
					Weight:       &r.opt.Weight,
					Priority:     &r.opt.Priority,
					Version:      &serviceInstance.Version,
					Metadata:     rmd,
					Healthy:      &r.opt.Healthy,
					Isolate:      &r.opt.Isolate,
					TTL:          &r.opt.TTL,
					Timeout:      &r.opt.Timeout,
					RetryCount:   &r.opt.RetryCount,
				},
			})
		if err != nil {
			return err
		}
		instanceID := service.InstanceID

		if r.opt.Heartbeat {
			// 启动心跳上报
			go func() {
				ticker := time.NewTicker(time.Second * time.Duration(r.opt.TTL))
				defer ticker.Stop()

				for {
					<-ticker.C

					err = r.provider.Heartbeat(&api.InstanceHeartbeatRequest{
						InstanceHeartbeatRequest: model.InstanceHeartbeatRequest{
							Service:      serviceInstance.Name + u.Scheme,
							Namespace:    r.opt.Namespace,
							Host:         host,
							Port:         portNum,
							ServiceToken: r.opt.ServiceToken,
							InstanceID:   instanceID,
							Timeout:      &r.opt.Timeout,
							RetryCount:   &r.opt.RetryCount,
						},
					})
					if err != nil {
						log.Error(err.Error())
						continue
					}
				}
			}()
		}

		ids = append(ids, instanceID)
	}
	// 设置 InstanceID 供注销使用
	serviceInstance.ID = strings.Join(ids, _instanceIDSeparator)
	return nil
}

// Deregister 注销服务
func (r *Registry) Deregister(_ context.Context, serviceInstance *registry.ServiceInstance) error {
	split := strings.Split(serviceInstance.ID, _instanceIDSeparator)
	for i, endpoint := range serviceInstance.Endpoints {
		// 解析 url
		u, err := url.Parse(endpoint)
		if err != nil {
			return err
		}

		// 解析 host 与 port
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return err
		}

		// 端口转为整数
		portNum, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
		// 调用注销
		err = r.provider.Deregister(
			&api.InstanceDeRegisterRequest{
				InstanceDeRegisterRequest: model.InstanceDeRegisterRequest{
					Service:      serviceInstance.Name + u.Scheme,
					ServiceToken: r.opt.ServiceToken,
					Namespace:    r.opt.Namespace,
					InstanceID:   split[i],
					Host:         host,
					Port:         portNum,
					Timeout:      &r.opt.Timeout,
					RetryCount:   &r.opt.RetryCount,
				},
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetService 按服务名获取实例列表
func (r *Registry) GetService(_ context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	// 获取全部实例
	instancesResponse, err := r.consumer.GetAllInstances(&api.GetAllInstancesRequest{
		GetAllInstancesRequest: model.GetAllInstancesRequest{
			Service:    serviceName,
			Namespace:  r.opt.Namespace,
			Timeout:    &r.opt.Timeout,
			RetryCount: &r.opt.RetryCount,
		},
	})
	if err != nil {
		return nil, err
	}

	serviceInstances := instancesToServiceInstances(instancesResponse.GetInstances())

	return serviceInstances, nil
}

// Watch 按服务名创建 watcher
func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return newWatcher(ctx, r.opt.Namespace, serviceName, r.consumer)
}

func instancesToServiceInstances(instances []model.Instance) []*registry.ServiceInstance {
	serviceInstances := make([]*registry.ServiceInstance, 0, len(instances))
	for _, instance := range instances {
		if instance.IsHealthy() {
			serviceInstances = append(serviceInstances, instanceToServiceInstance(instance))
		}
	}
	return serviceInstances
}

func instanceToServiceInstance(instance model.Instance) *registry.ServiceInstance {
	metadata := instance.GetMetadata()
	// 正常注册情况下不会出错
	kind := ""
	if k, ok := metadata["kind"]; ok {
		kind = k
	}
	return &registry.ServiceInstance{
		ID:        instance.GetId(),
		Name:      instance.GetService(),
		Version:   metadata["version"],
		Metadata:  metadata,
		Endpoints: []string{kind + "://" + net.JoinHostPort(instance.GetHost(), strconv.Itoa(int(instance.GetPort())))},
	}
}
