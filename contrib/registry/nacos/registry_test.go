package nacos

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/byteweap/wukong/component/registry"
)

// skipIfNacosNotAvailable 检查 Nacos 服务器是否可用，如果不可用则跳过测试
func skipIfNacosNotAvailable(t *testing.T) {
	t.Helper()

	// 从环境变量获取 Nacos 服务器地址，默认使用 localhost:8848
	nacosAddr := os.Getenv("NACOS_SERVER_ADDR")
	if nacosAddr == "" {
		nacosAddr = "127.0.0.1:8848"
	}

	// 尝试创建客户端连接
	reg, err := NewRegistry(ServerAddrs(nacosAddr))
	if err != nil {
		t.Skipf("Nacos server not available at %s: %v", nacosAddr, err)
	}
	reg.Close()

	// 尝试创建发现客户端
	discovery, err := NewDiscovery(ServerAddrs(nacosAddr))
	if err != nil {
		t.Skipf("Nacos server not available at %s: %v", nacosAddr, err)
	}
	discovery.Close()
}

// getNacosAddr 获取 Nacos 服务器地址
func getNacosAddr() string {
	addr := os.Getenv("NACOS_SERVER_ADDR")
	if addr == "" {
		return "127.0.0.1:8848"
	}
	return addr
}

func TestRegistry_Register(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	reg, err := NewRegistry(
		ServerAddrs(nacosAddr),
		Namespace("public"),
	)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	ctx := context.Background()
	service := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "test-service",
		Weight:    10,
		Version:   "v1.0.0",
		Metadata:  map[string]string{"env": "test"},
		Endpoints: []string{"http://127.0.0.1:8080"},
	}

	// 注册服务
	err = reg.Register(ctx, service)
	require.NoError(t, err)

	// 等待一下让注册生效
	time.Sleep(500 * time.Millisecond)

	// 验证服务已注册（通过 discovery 验证）
	discovery, err := NewDiscovery(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	instances, err := discovery.GetService(ctx, "test-service")
	require.NoError(t, err)
	require.Len(t, instances, 1)
	require.True(t, instances[0].Equal(service))

	// 清理：注销服务
	err = reg.Deregister(ctx, service)
	require.NoError(t, err)
}

func TestRegistry_Deregister(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	reg, err := NewRegistry(
		ServerAddrs(nacosAddr),
		Namespace("public"),
	)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	ctx := context.Background()
	service := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "test-service",
		Weight:    10,
		Version:   "v1.0.0",
		Endpoints: []string{"http://127.0.0.1:8080"},
	}

	// 注册服务
	err = reg.Register(ctx, service)
	require.NoError(t, err)

	// 等待一下让注册生效
	time.Sleep(500 * time.Millisecond)

	// 注销服务
	err = reg.Deregister(ctx, service)
	require.NoError(t, err)

	// 等待一下让注销生效
	time.Sleep(500 * time.Millisecond)

	// 验证服务已注销
	discovery, err := NewDiscovery(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	instances, err := discovery.GetService(ctx, "test-service")
	require.NoError(t, err)
	require.Len(t, instances, 0)
}

func TestRegistry_RegisterMultipleInstances(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	reg, err := NewRegistry(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	ctx := context.Background()

	// 注册多个实例
	services := []*registry.ServiceInstance{
		{ID: "instance-1", Name: "test-service", Weight: 10, Version: "v1.0.0", Endpoints: []string{"http://127.0.0.1:8080"}},
		{ID: "instance-2", Name: "test-service", Weight: 10, Version: "v1.0.0", Endpoints: []string{"http://127.0.0.1:8081"}},
		{ID: "instance-3", Name: "test-service", Weight: 10, Version: "v1.0.0", Endpoints: []string{"http://127.0.0.1:8082"}},
	}

	for _, service := range services {
		err = reg.Register(ctx, service)
		require.NoError(t, err)
	}

	// 等待一下让注册生效
	time.Sleep(500 * time.Millisecond)

	// 验证所有实例都已注册
	discovery, err := NewDiscovery(ServerAddrs(nacosAddr), Namespace("public"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	instances, err := discovery.GetService(ctx, "test-service")
	require.NoError(t, err)
	require.Len(t, instances, 3)

	// 清理：注销所有服务
	for _, service := range services {
		_ = reg.Deregister(ctx, service)
	}
}

func TestRegistry_Register_InvalidService(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	reg, err := NewRegistry(ServerAddrs(nacosAddr))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	ctx := context.Background()

	// 测试 nil service
	err = reg.Register(ctx, nil)
	require.Error(t, err)

	// 测试缺少 ID
	err = reg.Register(ctx, &registry.ServiceInstance{Name: "test"})
	require.Error(t, err)

	// 测试缺少 Name
	err = reg.Register(ctx, &registry.ServiceInstance{ID: "instance-1"})
	require.Error(t, err)

	// 测试缺少 Endpoints
	err = reg.Register(ctx, &registry.ServiceInstance{ID: "instance-1", Name: "test"})
	require.Error(t, err)
}

func TestRegistry_ID(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	reg, err := NewRegistry(ServerAddrs(nacosAddr))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	require.Equal(t, RegistryID, reg.ID())
}

func TestRegistry_Register_WithMetadata(t *testing.T) {
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

	// 注册服务
	err = reg.Register(ctx, service)
	require.NoError(t, err)

	// 等待一下让注册生效
	time.Sleep(500 * time.Millisecond)

	// 验证服务已注册且元数据正确
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
	err = reg.Deregister(ctx, service)
	require.NoError(t, err)
}

func TestRegistry_Close(t *testing.T) {
	skipIfNacosNotAvailable(t)

	nacosAddr := getNacosAddr()
	reg, err := NewRegistry(ServerAddrs(nacosAddr))
	require.NoError(t, err)

	ctx := context.Background()
	service := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "test-service",
		Endpoints: []string{"http://127.0.0.1:8080"},
	}

	// 注册服务
	err = reg.Register(ctx, service)
	require.NoError(t, err)

	// 关闭注册器
	err = reg.Close()
	require.NoError(t, err)

	// 再次关闭应该不会出错（幂等）
	err = reg.Close()
	require.NoError(t, err)

	// 清理：注销服务（使用新的注册器）
	reg2, _ := NewRegistry(ServerAddrs(nacosAddr))
	if reg2 != nil {
		_ = reg2.Deregister(ctx, service)
		_ = reg2.Close()
	}
}
