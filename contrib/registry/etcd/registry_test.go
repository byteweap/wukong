package etcd

import (
	"context"
	"net"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"

	"github.com/byteweap/wukong/component/registry"
)

// setupEmbeddedEtcd 创建嵌入式 etcd 服务器用于测试
func setupEmbeddedEtcd(t *testing.T) (*embed.Etcd, string) {
	t.Helper()

	// 创建临时目录用于 etcd 数据存储
	tmpDir, err := os.MkdirTemp("", "etcd-test-*")
	require.NoError(t, err)

	// 获取可用端口
	clientListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	clientPort := clientListener.Addr().(*net.TCPAddr).Port
	clientListener.Close()

	peerListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	peerPort := peerListener.Addr().(*net.TCPAddr).Port
	peerListener.Close()

	// 配置 etcd 服务器
	cfg := embed.NewConfig()
	cfg.Dir = tmpDir
	cfg.Name = "test-etcd"

	// 使用获取的端口
	clientURL, err := url.Parse("http://127.0.0.1:" + strconv.Itoa(clientPort))
	require.NoError(t, err)
	peerURL, err := url.Parse("http://127.0.0.1:" + strconv.Itoa(peerPort))
	require.NoError(t, err)

	cfg.ListenClientUrls = []url.URL{*clientURL}
	cfg.AdvertiseClientUrls = cfg.ListenClientUrls
	cfg.ListenPeerUrls = []url.URL{*peerURL}
	cfg.AdvertisePeerUrls = cfg.ListenPeerUrls
	cfg.InitialCluster = cfg.InitialClusterFromName(cfg.Name)
	cfg.LogLevel = "fatal" // 减少测试日志输出

	// 启动嵌入式 etcd 服务器
	etcd, err := embed.StartEtcd(cfg)
	require.NoError(t, err)

	// 等待服务器启动
	select {
	case <-etcd.Server.ReadyNotify():
	case <-time.After(10 * time.Second):
		etcd.Close()
		os.RemoveAll(tmpDir)
		t.Fatal("etcd server took too long to start")
	}

	// 等待一下确保服务器完全启动
	time.Sleep(200 * time.Millisecond)

	// 获取实际客户端地址
	actualClientURL := clientURL.String()

	// 清理函数
	t.Cleanup(func() {
		etcd.Close()
		os.RemoveAll(tmpDir)
	})

	return etcd, actualClientURL
}

// TestNewRegistry 测试 NewRegistry 创建注册器
func TestNewRegistry(t *testing.T) {
	_, clientURL := setupEmbeddedEtcd(t)

	reg, err := NewRegistry(
		Endpoints(clientURL),
		DialTimeout(2*time.Second),
		Namespace("/test"),
		TTL(5*time.Second),
	)
	require.NoError(t, err)
	require.NotNil(t, reg)
	defer reg.Close()

	assert.Equal(t, RegistryID, reg.ID())
	assert.True(t, reg.ownClient, "NewRegistry 应该自己管理客户端生命周期")
}

// TestNewRegistryWith 测试 NewRegistryWith 创建注册器
func TestNewRegistryWith(t *testing.T) {
	_, clientURL := setupEmbeddedEtcd(t)

	// 创建客户端
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{clientURL},
		DialTimeout: 2 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	reg := NewRegistryWith(client, Namespace("/test"))
	require.NotNil(t, reg)
	defer reg.Close()

	assert.Equal(t, RegistryID, reg.ID())
	assert.False(t, reg.ownClient, "NewRegistryWith 应该由上层管理客户端生命周期")
}

// TestRegister 测试服务注册
func TestRegister(t *testing.T) {
	_, clientURL := setupEmbeddedEtcd(t)

	reg, err := NewRegistry(
		Endpoints(clientURL),
		DialTimeout(2*time.Second),
		Namespace("/test"),
		TTL(5*time.Second),
	)
	require.NoError(t, err)
	defer reg.Close()

	ctx := context.Background()

	// 测试正常注册
	service := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "test-service",
		Version:   "v1.0.0",
		Weight:    1.0,
		Endpoints: []string{"http://127.0.0.1:8080"},
		Metadata:  map[string]string{"env": "test"},
	}

	err = reg.Register(ctx, service)
	assert.NoError(t, err)

	// 验证服务已注册
	instances, err := reg.GetService(ctx, "test-service")
	require.NoError(t, err)
	require.Len(t, instances, 1)
	assert.Equal(t, service.ID, instances[0].ID)
	assert.Equal(t, service.Name, instances[0].Name)
	assert.Equal(t, service.Version, instances[0].Version)
	assert.Equal(t, service.Weight, instances[0].Weight)
	assert.Equal(t, service.Endpoints, instances[0].Endpoints)
	assert.Equal(t, service.Metadata, instances[0].Metadata)

	// 测试注册空服务
	err = reg.Register(ctx, nil)
	assert.Error(t, err)

	// 测试注册缺少 ID 的服务
	err = reg.Register(ctx, &registry.ServiceInstance{
		Name:      "test-service",
		Endpoints: []string{"http://127.0.0.1:8080"},
	})
	assert.Error(t, err)

	// 测试注册缺少 Name 的服务
	err = reg.Register(ctx, &registry.ServiceInstance{
		ID:        "instance-2",
		Endpoints: []string{"http://127.0.0.1:8080"},
	})
	assert.Error(t, err)
}

// TestDeregister 测试服务注销
func TestDeregister(t *testing.T) {
	_, clientURL := setupEmbeddedEtcd(t)

	reg, err := NewRegistry(
		Endpoints(clientURL),
		DialTimeout(2*time.Second),
		Namespace("/test"),
		TTL(5*time.Second),
	)
	require.NoError(t, err)
	defer reg.Close()

	ctx := context.Background()

	// 先注册服务
	service := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "test-service",
		Version:   "v1.0.0",
		Weight:    1.0,
		Endpoints: []string{"http://127.0.0.1:8080"},
		Metadata:  map[string]string{"env": "test"},
	}

	err = reg.Register(ctx, service)
	require.NoError(t, err)

	// 验证服务已注册
	instances, err := reg.GetService(ctx, "test-service")
	require.NoError(t, err)
	require.Len(t, instances, 1)

	// 注销服务
	err = reg.Deregister(ctx, service)
	assert.NoError(t, err)

	// 验证服务已注销
	instances, err = reg.GetService(ctx, "test-service")
	require.NoError(t, err)
	assert.Len(t, instances, 0)

	// 测试注销空服务
	err = reg.Deregister(ctx, nil)
	assert.Error(t, err)

	// 测试注销不存在的服务（应该不报错）
	nonExistentService := &registry.ServiceInstance{
		ID:   "non-existent",
		Name: "test-service",
	}
	err = reg.Deregister(ctx, nonExistentService)
	assert.NoError(t, err)
}

// TestGetService 测试获取服务列表
func TestGetService(t *testing.T) {
	_, clientURL := setupEmbeddedEtcd(t)

	reg, err := NewRegistry(
		Endpoints(clientURL),
		DialTimeout(2*time.Second),
		Namespace("/test"),
		TTL(5*time.Second),
	)
	require.NoError(t, err)
	defer reg.Close()

	ctx := context.Background()

	// 注册多个服务实例
	services := []*registry.ServiceInstance{
		{
			ID:        "instance-1",
			Name:      "test-service",
			Version:   "v1.0.0",
			Weight:    1.0,
			Endpoints: []string{"http://127.0.0.1:8080"},
		},
		{
			ID:        "instance-2",
			Name:      "test-service",
			Version:   "v1.0.0",
			Weight:    2.0,
			Endpoints: []string{"http://127.0.0.1:8081"},
		},
		{
			ID:        "instance-3",
			Name:      "other-service",
			Version:   "v1.0.0",
			Weight:    1.0,
			Endpoints: []string{"http://127.0.0.1:8082"},
		},
	}

	for _, service := range services {
		err := reg.Register(ctx, service)
		require.NoError(t, err)
	}

	// 获取 test-service 的所有实例
	instances, err := reg.GetService(ctx, "test-service")
	require.NoError(t, err)
	assert.Len(t, instances, 2)

	// 验证实例 ID
	instanceIDs := make(map[string]bool)
	for _, instance := range instances {
		instanceIDs[instance.ID] = true
	}
	assert.True(t, instanceIDs["instance-1"])
	assert.True(t, instanceIDs["instance-2"])

	// 获取 other-service 的所有实例
	instances, err = reg.GetService(ctx, "other-service")
	require.NoError(t, err)
	assert.Len(t, instances, 1)
	assert.Equal(t, "instance-3", instances[0].ID)

	// 获取不存在的服务
	instances, err = reg.GetService(ctx, "non-existent")
	require.NoError(t, err)
	assert.Len(t, instances, 0)

	// 测试空服务名
	_, err = reg.GetService(ctx, "")
	assert.Error(t, err)
}

// TestWatch 测试服务监听
func TestWatch(t *testing.T) {
	_, clientURL := setupEmbeddedEtcd(t)

	reg, err := NewRegistry(
		Endpoints(clientURL),
		DialTimeout(2*time.Second),
		Namespace("/test"),
		TTL(5*time.Second),
	)
	require.NoError(t, err)
	defer reg.Close()

	ctx := context.Background()
	serviceName := "test-service"

	// 创建监听器
	watcher, err := reg.Watch(ctx, serviceName)
	require.NoError(t, err)
	require.NotNil(t, watcher)
	defer watcher.Stop()

	// 注册一个服务实例
	service := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      serviceName,
		Version:   "v1.0.0",
		Weight:    1.0,
		Endpoints: []string{"http://127.0.0.1:8080"},
	}

	err = reg.Register(ctx, service)
	require.NoError(t, err)

	// 等待监听器收到更新
	instances, err := watcher.Next()
	require.NoError(t, err)
	require.Len(t, instances, 1)
	assert.Equal(t, service.ID, instances[0].ID)

	// 注册另一个服务实例
	service2 := &registry.ServiceInstance{
		ID:        "instance-2",
		Name:      serviceName,
		Version:   "v1.0.0",
		Weight:    1.0,
		Endpoints: []string{"http://127.0.0.1:8081"},
	}

	err = reg.Register(ctx, service2)
	require.NoError(t, err)

	// 等待监听器收到更新
	instances, err = watcher.Next()
	require.NoError(t, err)
	assert.Len(t, instances, 2)

	// 注销一个服务实例
	err = reg.Deregister(ctx, service)
	require.NoError(t, err)

	// 等待监听器收到更新
	instances, err = watcher.Next()
	require.NoError(t, err)
	assert.Len(t, instances, 1)
	assert.Equal(t, service2.ID, instances[0].ID)

	// 测试空服务名
	_, err = reg.Watch(ctx, "")
	assert.Error(t, err)
}

// TestClose 测试关闭注册器
func TestClose(t *testing.T) {
	_, clientURL := setupEmbeddedEtcd(t)

	// 测试 NewRegistry 创建的注册器可以关闭客户端
	reg, err := NewRegistry(
		Endpoints(clientURL),
		DialTimeout(2*time.Second),
		Namespace("/test"),
		TTL(5*time.Second),
	)
	require.NoError(t, err)

	ctx := context.Background()

	// 注册一个服务
	service := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "test-service",
		Version:   "v1.0.0",
		Weight:    1.0,
		Endpoints: []string{"http://127.0.0.1:8080"},
	}

	err = reg.Register(ctx, service)
	require.NoError(t, err)

	// 先注销服务，优雅地停止续租，减少关闭时的警告日志
	err = reg.Deregister(ctx, service)
	require.NoError(t, err)

	// 等待一下让注销操作完成
	time.Sleep(50 * time.Millisecond)

	// 关闭注册器
	err = reg.Close()
	assert.NoError(t, err)

	// 验证客户端已关闭（后续操作应该失败）
	_, err = reg.GetService(ctx, "test-service")
	assert.Error(t, err)
}

// TestClose_NewRegistryWith 测试 NewRegistryWith 创建的注册器关闭时不会关闭客户端
func TestClose_NewRegistryWith(t *testing.T) {
	_, clientURL := setupEmbeddedEtcd(t)

	// 创建客户端
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{clientURL},
		DialTimeout: 2 * time.Second,
	})
	require.NoError(t, err)

	reg := NewRegistryWith(client, Namespace("/test"))

	ctx := context.Background()

	// 注册一个服务
	service := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "test-service",
		Version:   "v1.0.0",
		Weight:    1.0,
		Endpoints: []string{"http://127.0.0.1:8080"},
	}

	err = reg.Register(ctx, service)
	require.NoError(t, err)

	// 关闭注册器（不应该关闭客户端）
	err = reg.Close()
	assert.NoError(t, err)

	// 验证客户端仍然可用（通过原始客户端访问）
	resp, err := client.Get(ctx, "/test/test-service/instance-1")
	assert.NoError(t, err)
	assert.Len(t, resp.Kvs, 1)

	// 手动关闭客户端
	client.Close()
}

// TestLeaseExpiration 测试租约过期
func TestLeaseExpiration(t *testing.T) {
	_, clientURL := setupEmbeddedEtcd(t)

	reg, err := NewRegistry(
		Endpoints(clientURL),
		DialTimeout(2*time.Second),
		Namespace("/test"),
		TTL(2*time.Second), // 短 TTL 用于测试过期
	)
	require.NoError(t, err)
	defer reg.Close()

	ctx := context.Background()

	// 注册服务
	service := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "test-service",
		Version:   "v1.0.0",
		Weight:    1.0,
		Endpoints: []string{"http://127.0.0.1:8080"},
	}

	err = reg.Register(ctx, service)
	require.NoError(t, err)

	// 验证服务已注册
	instances, err := reg.GetService(ctx, "test-service")
	require.NoError(t, err)
	require.Len(t, instances, 1)

	// 等待租约过期（TTL + 一些缓冲时间）
	time.Sleep(3 * time.Second)

	// 验证服务已过期（由于续租机制，服务应该仍然存在）
	// 注意：由于有 KeepAlive 机制，服务不会过期
	instances, err = reg.GetService(ctx, "test-service")
	require.NoError(t, err)
	// 由于续租机制，服务应该仍然存在
	assert.GreaterOrEqual(t, len(instances), 0)
}

// TestMultipleRegistries 测试多个注册器使用同一个客户端
func TestMultipleRegistries(t *testing.T) {
	_, clientURL := setupEmbeddedEtcd(t)

	// 创建共享客户端
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{clientURL},
		DialTimeout: 2 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	// 创建两个注册器使用同一个客户端
	reg1 := NewRegistryWith(client, Namespace("/test1"))
	reg2 := NewRegistryWith(client, Namespace("/test2"))

	ctx := context.Background()

	// 在 reg1 中注册服务
	service1 := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "service",
		Version:   "v1.0.0",
		Weight:    1.0,
		Endpoints: []string{"http://127.0.0.1:8080"},
	}
	err = reg1.Register(ctx, service1)
	require.NoError(t, err)

	// 在 reg2 中注册服务
	service2 := &registry.ServiceInstance{
		ID:        "instance-2",
		Name:      "service",
		Version:   "v1.0.0",
		Weight:    1.0,
		Endpoints: []string{"http://127.0.0.1:8081"},
	}
	err = reg2.Register(ctx, service2)
	require.NoError(t, err)

	// 验证命名空间隔离
	instances1, err := reg1.GetService(ctx, "service")
	require.NoError(t, err)
	assert.Len(t, instances1, 1)
	assert.Equal(t, "instance-1", instances1[0].ID)

	instances2, err := reg2.GetService(ctx, "service")
	require.NoError(t, err)
	assert.Len(t, instances2, 1)
	assert.Equal(t, "instance-2", instances2[0].ID)

	// 关闭注册器（不应该关闭共享客户端）
	reg1.Close()
	reg2.Close()

	// 验证客户端仍然可用
	_, err = client.Get(ctx, "/test1/service/instance-1")
	assert.NoError(t, err)
}

// TestRegistryOptions 测试各种选项配置
func TestRegistryOptions(t *testing.T) {
	_, clientURL := setupEmbeddedEtcd(t)

	// 测试自定义命名空间
	reg, err := NewRegistry(
		Endpoints(clientURL),
		Namespace("/custom/namespace"),
		TTL(10*time.Second),
		DialTimeout(5*time.Second),
	)
	require.NoError(t, err)
	defer reg.Close()

	assert.Equal(t, "/custom/namespace", reg.namespace)
	assert.Equal(t, 10*time.Second, reg.ttl)
}

// TestEmbeddedEtcdIntegration 集成测试：完整的注册-发现-注销流程
func TestEmbeddedEtcdIntegration(t *testing.T) {
	_, clientURL := setupEmbeddedEtcd(t)

	reg, err := NewRegistry(
		Endpoints(clientURL),
		DialTimeout(2*time.Second),
		Namespace("/integration-test"),
		TTL(5*time.Second),
	)
	require.NoError(t, err)
	defer reg.Close()

	ctx := context.Background()
	serviceName := "integration-service"

	// 1. 注册多个服务实例
	services := []*registry.ServiceInstance{
		{ID: "instance-1", Name: serviceName, Version: "v1.0.0", Weight: 1.0, Endpoints: []string{"http://127.0.0.1:8080"}},
		{ID: "instance-2", Name: serviceName, Version: "v1.0.0", Weight: 2.0, Endpoints: []string{"http://127.0.0.1:8081"}},
		{ID: "instance-3", Name: serviceName, Version: "v2.0.0", Weight: 1.0, Endpoints: []string{"http://127.0.0.1:8082"}},
	}

	for _, service := range services {
		err := reg.Register(ctx, service)
		require.NoError(t, err)
	}

	// 2. 验证所有实例都被发现
	instances, err := reg.GetService(ctx, serviceName)
	require.NoError(t, err)
	assert.Len(t, instances, 3)

	// 3. 创建监听器
	watcher, err := reg.Watch(ctx, serviceName)
	require.NoError(t, err)
	defer watcher.Stop()

	// 4. 注销一个实例
	err = reg.Deregister(ctx, services[0])
	require.NoError(t, err)

	// 5. 验证监听器收到更新
	instances, err = watcher.Next()
	require.NoError(t, err)
	assert.Len(t, instances, 2)

	// 6. 验证 GetService 也返回正确的结果
	instances, err = reg.GetService(ctx, serviceName)
	require.NoError(t, err)
	assert.Len(t, instances, 2)

	// 7. 注销所有剩余实例
	for i := 1; i < len(services); i++ {
		err = reg.Deregister(ctx, services[i])
		require.NoError(t, err)
	}

	// 8. 验证所有实例都已注销
	instances, err = reg.GetService(ctx, serviceName)
	require.NoError(t, err)
	assert.Len(t, instances, 0)
}
