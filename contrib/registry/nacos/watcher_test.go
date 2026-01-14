package nacos

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/byteweap/wukong/component/registry"
)

func TestWatcher_Next_InitialLoad(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	reg, err := NewRegistry(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	discovery, err := NewDiscovery(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	ctx := context.Background()

	// 先注册一些服务
	services := []*registry.ServiceInstance{
		{ID: "instance-1", Name: "test-service", Weight: 10, Endpoints: []string{"http://127.0.0.1:8080"}},
		{ID: "instance-2", Name: "test-service", Weight: 10, Endpoints: []string{"http://127.0.0.1:8081"}},
	}

	for _, service := range services {
		err = reg.Register(ctx, service)
		require.NoError(t, err)
	}

	// 等待一下让注册生效
	time.Sleep(500 * time.Millisecond)

	// 创建监听器
	watcher, err := discovery.Watch(ctx, "test-service")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, watcher.Close()) })

	// 首次 Next 应该立即返回当前实例列表
	instances, err := watcher.Next()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(instances), 2)

	// 清理
	for _, service := range services {
		_ = reg.Deregister(ctx, service)
	}
}

func TestWatcher_Next_BlockUntilChange(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	reg, err := NewRegistry(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	discovery, err := NewDiscovery(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	ctx := context.Background()

	// 创建监听器（此时可能没有服务）
	watcher, err := discovery.Watch(ctx, "test-service")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, watcher.Close()) })

	// 在另一个 goroutine 中注册服务
	service := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "test-service",
		Endpoints: []string{"http://127.0.0.1:8080"},
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		err := reg.Register(ctx, service)
		require.NoError(t, err)
	}()

	// Next 应该阻塞直到服务注册
	instances, err := watcher.Next()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(instances), 1)
	found := false
	for _, inst := range instances {
		if inst.ID == service.ID {
			require.True(t, inst.Equal(service))
			found = true
			break
		}
	}
	require.True(t, found, "应该找到注册的服务实例")

	// 清理
	_ = reg.Deregister(ctx, service)
}

func TestWatcher_Next_ServiceUpdate(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	reg, err := NewRegistry(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	discovery, err := NewDiscovery(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	ctx := context.Background()

	// 先注册服务
	service := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "test-service",
		Endpoints: []string{"http://127.0.0.1:8080"},
	}
	err = reg.Register(ctx, service)
	require.NoError(t, err)

	// 等待一下让注册生效
	time.Sleep(500 * time.Millisecond)

	// 创建监听器
	watcher, err := discovery.Watch(ctx, "test-service")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, watcher.Close()) })

	// 首次获取
	instances, err := watcher.Next()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(instances), 1)

	// 更新服务（重新注册，但 endpoints 不同）
	updatedService := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "test-service",
		Endpoints: []string{"http://127.0.0.1:8080", "grpc://127.0.0.1:9000"},
	}
	err = reg.Deregister(ctx, service)
	require.NoError(t, err)
	time.Sleep(200 * time.Millisecond)
	err = reg.Register(ctx, updatedService)
	require.NoError(t, err)

	// Next 应该返回更新后的服务
	instances, err = watcher.Next()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(instances), 1)
	found := false
	for _, inst := range instances {
		if inst.ID == updatedService.ID {
			require.True(t, inst.Equal(updatedService))
			found = true
			break
		}
	}
	require.True(t, found, "应该找到更新后的服务实例")

	// 清理
	_ = reg.Deregister(ctx, updatedService)
}

func TestWatcher_Next_ServiceDelete(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	reg, err := NewRegistry(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	discovery, err := NewDiscovery(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	ctx := context.Background()

	// 注册两个服务
	service1 := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "test-service",
		Endpoints: []string{"http://127.0.0.1:8080"},
	}
	service2 := &registry.ServiceInstance{
		ID:        "instance-2",
		Name:      "test-service",
		Endpoints: []string{"http://127.0.0.1:8081"},
	}

	err = reg.Register(ctx, service1)
	require.NoError(t, err)
	err = reg.Register(ctx, service2)
	require.NoError(t, err)

	// 等待一下让注册生效
	time.Sleep(500 * time.Millisecond)

	// 创建监听器
	watcher, err := discovery.Watch(ctx, "test-service")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, watcher.Close()) })

	// 首次获取
	instances, err := watcher.Next()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(instances), 2)

	// 删除一个服务
	err = reg.Deregister(ctx, service1)
	require.NoError(t, err)

	// Next 应该返回更新后的列表
	instances, err = watcher.Next()
	require.NoError(t, err)
	// 应该只剩下 service2
	found := false
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

func TestWatcher_Close(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	discovery, err := NewDiscovery(ServerAddrs(nacosAddr))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	ctx := context.Background()

	watcher, err := discovery.Watch(ctx, "test-service")
	require.NoError(t, err)

	// 停止监听器
	err = watcher.Close()
	require.NoError(t, err)

	// 再次停止应该不会出错（幂等）
	err = watcher.Close()
	require.NoError(t, err)

	// 停止后 Next 应该返回错误
	_, err = watcher.Next()
	require.Error(t, err)
}

func TestWatcher_Close_WithContextCancel(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	discovery, err := NewDiscovery(ServerAddrs(nacosAddr))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	ctx, cancel := context.WithCancel(context.Background())

	watcher, err := discovery.Watch(ctx, "test-service")
	require.NoError(t, err)

	// 取消 context
	cancel()

	// Next 应该返回 context 取消错误
	_, err = watcher.Next()
	require.Error(t, err)
	require.Equal(t, context.Canceled, err)

	// Close 应该正常工作
	err = watcher.Close()
	require.NoError(t, err)
}

func TestWatcher_ID(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	discovery, err := NewDiscovery(ServerAddrs(nacosAddr))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	ctx := context.Background()
	watcher, err := discovery.Watch(ctx, "test-service")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, watcher.Close()) })

	require.Equal(t, WatcherID, watcher.ID())
}

func TestWatcher_ConcurrentOperations(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	reg, err := NewRegistry(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	discovery, err := NewDiscovery(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	ctx := context.Background()

	watcher, err := discovery.Watch(ctx, "test-service")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, watcher.Close()) })

	// 并发注册多个服务
	const numServices = 10
	services := make([]*registry.ServiceInstance, numServices)
	for i := 0; i < numServices; i++ {
		services[i] = &registry.ServiceInstance{
			ID:        fmt.Sprintf("instance-%d", i),
			Name:      "test-service",
			Endpoints: []string{fmt.Sprintf("http://127.0.0.1:%d", 8080+i)},
		}
		go func(s *registry.ServiceInstance) {
			_ = reg.Register(ctx, s)
		}(services[i])
	}

	// 等待所有服务注册
	time.Sleep(1 * time.Second)

	// Next 应该返回所有服务
	instances, err := watcher.Next()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(instances), 1) // 至少有一个服务

	// 清理
	for _, service := range services {
		_ = reg.Deregister(ctx, service)
	}
}
