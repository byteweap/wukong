package nacos

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/byteweap/wukong/component/registry"
)

func TestDiscovery_GetService(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	// 先注册一些服务
	reg, err := NewRegistry(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	ctx := context.Background()
	services := []*registry.ServiceInstance{
		{ID: "instance-1", Name: "service-a", Weight: 10, Version: "v1.0.0", Endpoints: []string{"http://127.0.0.1:8080"}},
		{ID: "instance-2", Name: "service-a", Weight: 10, Version: "v1.0.0", Endpoints: []string{"http://127.0.0.1:8081"}},
		{ID: "instance-3", Name: "service-b", Weight: 10, Version: "v1.0.0", Endpoints: []string{"http://127.0.0.1:8082"}},
	}

	for _, service := range services {
		err = reg.Register(ctx, service)
		require.NoError(t, err)
	}

	// 等待一下让注册生效
	time.Sleep(500 * time.Millisecond)

	// 创建发现器
	discovery, err := NewDiscovery(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	// 获取 service-a 的所有实例
	instances, err := discovery.GetService(ctx, "service-a")
	require.NoError(t, err)
	require.Len(t, instances, 2)
	for _, instance := range instances {
		t.Logf("service-a ----- instance: %+v", instance)
	}

	// 获取 service-b 的所有实例
	instances, err = discovery.GetService(ctx, "service-b")
	require.NoError(t, err)
	require.Len(t, instances, 1)
	for _, instance := range instances {
		t.Logf("service-b ----- instance: %+v", instance)
	}

	// 获取不存在的服务
	instances, err = discovery.GetService(ctx, "service-c")
	require.NoError(t, err)
	require.Len(t, instances, 0)

	// 清理
	for _, service := range services {
		_ = reg.Deregister(ctx, service)
	}
}

func TestDiscovery_GetService_EmptyName(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	discovery, err := NewDiscovery(ServerAddrs(nacosAddr))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	ctx := context.Background()
	_, err = discovery.GetService(ctx, "")
	require.Error(t, err)
}

func TestDiscovery_Watch(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	reg, err := NewRegistry(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	discovery, err := NewDiscovery(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	ctx := context.Background()

	// 先注册一个服务
	service1 := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "test-service",
		Endpoints: []string{"http://127.0.0.1:8080"},
	}
	err = reg.Register(ctx, service1)
	require.NoError(t, err)

	// 等待一下让注册生效
	time.Sleep(500 * time.Millisecond)

	// 创建监听器
	watcher, err := discovery.Watch(ctx, "test-service")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, watcher.Close()) })

	// Next 应该返回当前实例
	instances, err := watcher.Next()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(instances), 1) // 至少有一个实例
	found := false
	for _, inst := range instances {
		if inst.ID == service1.ID {
			require.True(t, inst.Equal(service1))
			found = true
			break
		}
	}
	require.True(t, found, "应该找到注册的服务实例")

	// 注册新实例
	service2 := &registry.ServiceInstance{
		ID:        "instance-2",
		Name:      "test-service",
		Endpoints: []string{"http://127.0.0.1:8081"},
	}
	err = reg.Register(ctx, service2)
	require.NoError(t, err)

	// Next 应该返回更新后的实例列表
	instances, err = watcher.Next()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(instances), 2) // 至少有两个实例

	// 注销一个实例
	err = reg.Deregister(ctx, service1)
	require.NoError(t, err)

	// Next 应该返回更新后的实例列表
	instances, err = watcher.Next()
	require.NoError(t, err)
	// 应该只剩下 service2
	found = false
	for _, inst := range instances {
		if inst.ID == service2.ID {
			require.True(t, inst.Equal(service2))
			found = true
			break
		}
	}
	require.True(t, found, "应该找到剩余的服务实例")

	// 清理
	_ = reg.Deregister(ctx, service2)
}

func TestDiscovery_Watch_EmptyName(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	discovery, err := NewDiscovery(ServerAddrs(nacosAddr))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	ctx := context.Background()
	_, err = discovery.Watch(ctx, "")
	require.Error(t, err)
}

func TestDiscovery_ID(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	discovery, err := NewDiscovery(ServerAddrs(nacosAddr))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	require.Equal(t, DiscoveryID, discovery.ID())
}

func TestDiscovery_GetService_WithMetadata(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	reg, err := NewRegistry(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	ctx := context.Background()
	service := &registry.ServiceInstance{
		ID:      "instance-1",
		Name:    "test-service",
		Version: "v1.0.0",
		Metadata: map[string]string{
			"env":      "production",
			"region":   "us-east-1",
			"instance": "web-01",
		},
		Endpoints: []string{"http://127.0.0.1:8080", "grpc://127.0.0.1:9000"},
	}

	err = reg.Register(ctx, service)
	require.NoError(t, err)

	// 等待一下让注册生效
	time.Sleep(500 * time.Millisecond)

	discovery, err := NewDiscovery(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	instances, err := discovery.GetService(ctx, "test-service")
	require.NoError(t, err)
	require.Len(t, instances, 1)
	require.True(t, instances[0].Equal(service))
	require.Equal(t, service.Metadata, instances[0].Metadata)
	require.Equal(t, service.Endpoints, instances[0].Endpoints)

	// 清理
	_ = reg.Deregister(ctx, service)
}
