package servicecomb

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/component/registry"
	"github.com/go-chassis/cari/discovery"
	"github.com/go-chassis/cari/pkg/errsvc"
	"github.com/go-chassis/sc-client"
	"github.com/gofrs/uuid"
)

func init() {
	appID = os.Getenv(appIDVar)
	if appID == "" {
		appID = "default"
	}
	env = os.Getenv(envVar)
}

var (
	curServiceID string
	appID        string
	env          string
)

const (
	appIDKey         = "appId"
	envKey           = "environment"
	envVar           = "CAS_ENVIRONMENT_ID"
	appIDVar         = "CAS_APPLICATION_NAME"
	frameWorkName    = "wukong"
	frameWorkVersion = "v1"
)

// Registry 是 servicecomb 注册中心实现
type Registry struct {
	cli RegistryClient
}

var _ registry.Registry = (*Registry)(nil)

func NewRegistry(client RegistryClient) *Registry {
	r := &Registry{
		cli: client,
	}
	return r
}

func (r *Registry) ID() string {
	return "servicecomb"
}

func (r *Registry) GetService(_ context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	instances, err := r.cli.FindMicroServiceInstances("", appID, serviceName, "")
	if err != nil {
		return nil, err
	}
	svcInstances := make([]*registry.ServiceInstance, 0, len(instances))
	for _, instance := range instances {
		svcInstances = append(svcInstances, &registry.ServiceInstance{
			ID:        instance.InstanceId,
			Name:      serviceName,
			Metadata:  instance.Properties,
			Endpoints: instance.Endpoints,
			Version:   instance.ServiceId,
		})
	}
	return svcInstances, nil
}

func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return newWatcher(ctx, r.cli, serviceName)
}

func (r *Registry) Register(_ context.Context, svcIns *registry.ServiceInstance) error {
	fw := &discovery.FrameWork{
		Name:    frameWorkName,
		Version: frameWorkVersion,
	}
	ms := &discovery.MicroService{
		ServiceName: svcIns.Name,
		AppId:       appID,
		Version:     svcIns.Version,
		Environment: env,
		Framework:   fw,
	}
	// 尝试注册微服务
	sid, err := r.cli.RegisterService(ms)
	// 失败时可能是服务已存在
	if err != nil {
		registryException, ok := err.(*sc.RegistryException)
		if !ok {
			return err
		}
		var svcErr errsvc.Error
		parseErr := json.Unmarshal([]byte(registryException.Message), &svcErr)
		if parseErr != nil {
			return parseErr
		}
		// 错误码不是服务已存在时直接返回
		if svcErr.Code != discovery.ErrServiceAlreadyExists {
			return err
		}
		sid, err = r.cli.GetMicroServiceID(appID, ms.ServiceName, ms.Version, ms.Environment)
		if err != nil {
			return err
		}
	} else {
		// 记录新注册服务的 ID
		curServiceID = sid
	}
	if svcIns.ID == "" {
		var id uuid.UUID
		id, err = uuid.NewV4()
		if err != nil {
			return err
		}
		svcIns.ID = id.String()
	}
	props := map[string]string{
		appIDKey: appID,
		envKey:   env,
	}
	_, err = r.cli.RegisterMicroServiceInstance(&discovery.MicroServiceInstance{
		InstanceId: svcIns.ID,
		ServiceId:  sid,
		Endpoints:  svcIns.Endpoints,
		HostName:   svcIns.ID,
		Properties: props,
		Version:    svcIns.Version,
	})
	if err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			_, err = r.cli.Heartbeat(sid, svcIns.ID)
			if err != nil {
				log.Errorf("failed to send heartbeat: %v", err)
				continue
			}
		}
	}()
	return nil
}

func (r *Registry) Deregister(_ context.Context, svcIns *registry.ServiceInstance) error {
	sid, err := r.cli.GetMicroServiceID(appID, svcIns.Name, svcIns.Version, env)
	if err != nil {
		return err
	}
	_, err = r.cli.UnregisterMicroServiceInstance(sid, svcIns.ID)
	if err != nil {
		return err
	}
	return nil
}
