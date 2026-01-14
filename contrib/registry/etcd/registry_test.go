package etcd

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/server/v3/embed"

	"github.com/byteweap/wukong/component/registry"
)

// runEtcdServer 启动一个嵌入式的 etcd 服务器用于测试
func runEtcdServer(t *testing.T) (*embed.Etcd, string) {
	t.Helper()

	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "etcd-test-*")
	require.NoError(t, err)

	cfg := embed.NewConfig()
	cfg.Dir = filepath.Join(tmpDir, "etcd")
	cfg.LogLevel = "fatal"              // 只输出致命错误，减少测试日志输出
	cfg.LogOutputs = []string{"stderr"} // Windows 不支持 /dev/null
	cfg.Name = "test-etcd"

	// 使用随机端口
	clientURL, _ := url.Parse("http://127.0.0.1:0")
	peerURL, _ := url.Parse("http://127.0.0.1:0")
	cfg.ListenClientUrls = []url.URL{*clientURL}
	cfg.AdvertiseClientUrls = []url.URL{*clientURL}
	cfg.ListenPeerUrls = []url.URL{*peerURL}
	cfg.AdvertisePeerUrls = []url.URL{*peerURL}

	// 设置初始集群配置
	cfg.InitialCluster = cfg.InitialClusterFromName(cfg.Name)
	cfg.InitialClusterToken = "test-cluster-token"

	e, err := embed.StartEtcd(cfg)
	require.NoError(t, err)

	// 等待服务器就绪
	select {
	case <-e.Server.ReadyNotify():
		// 服务器已就绪
	case <-time.After(10 * time.Second):
		e.Close()
		t.Fatal("etcd server failed to start within 10 seconds")
	}

	// 获取客户端地址（使用实际监听的地址）
	clientAddr := e.Clients[0].Addr().String()
	clientURLStr := "http://" + clientAddr

	t.Cleanup(func() {
		// 优雅关闭 etcd 服务器
		// 先停止接受新连接
		e.Server.Stop()

		// 等待服务器完全停止
		select {
		case <-e.Server.StopNotify():
			// 服务器已停止
		case <-time.After(2 * time.Second):
			// 超时，继续关闭
		}

		// 关闭 etcd 实例（这会关闭所有连接）
		e.Close()

		// 等待一下，让所有 goroutine 完成
		time.Sleep(100 * time.Millisecond)

		// 清理临时目录
		os.RemoveAll(tmpDir)
	})

	return e, clientURLStr
}

func TestRegistry_Register(t *testing.T) {
	_, clientURL := runEtcdServer(t)

	reg, err := NewRegistry(
		Endpoints(clientURL),
		TTL(5*time.Second),
	)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	ctx := context.Background()
	service := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "test-service",
		Version:   "v1.0.0",
		Metadata:  map[string]string{"env": "test"},
		Endpoints: []string{"http://127.0.0.1:8080"},
	}

	// 注册服务
	err = reg.Register(ctx, service)
	require.NoError(t, err)

	// 验证服务已注册（通过 discovery 验证）
	discovery, err := NewDiscovery(Endpoints(clientURL))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	instances, err := discovery.GetService(ctx, "test-service")
	require.NoError(t, err)
	require.Len(t, instances, 1)
	require.True(t, instances[0].Equal(service))
}

func TestRegistry_Deregister(t *testing.T) {
	_, clientURL := runEtcdServer(t)

	reg, err := NewRegistry(
		Endpoints(clientURL),
		TTL(5*time.Second),
	)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	ctx := context.Background()
	service := &registry.ServiceInstance{
		ID:        "instance-1",
		Name:      "test-service",
		Version:   "v1.0.0",
		Endpoints: []string{"http://127.0.0.1:8080"},
	}

	// 注册服务
	err = reg.Register(ctx, service)
	require.NoError(t, err)

	// 注销服务
	err = reg.Deregister(ctx, service)
	require.NoError(t, err)

	// 验证服务已注销
	discovery, err := NewDiscovery(Endpoints(clientURL))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	instances, err := discovery.GetService(ctx, "test-service")
	require.NoError(t, err)
	require.Len(t, instances, 0)
}

func TestRegistry_RegisterMultipleInstances(t *testing.T) {
	_, clientURL := runEtcdServer(t)

	reg, err := NewRegistry(Endpoints(clientURL))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	ctx := context.Background()

	// 注册多个实例
	services := []*registry.ServiceInstance{
		{ID: "instance-1", Name: "test-service", Version: "v1.0.0", Endpoints: []string{"http://127.0.0.1:8080"}},
		{ID: "instance-2", Name: "test-service", Version: "v1.0.0", Endpoints: []string{"http://127.0.0.1:8081"}},
		{ID: "instance-3", Name: "test-service", Version: "v1.0.0", Endpoints: []string{"http://127.0.0.1:8082"}},
	}

	for _, service := range services {
		err = reg.Register(ctx, service)
		require.NoError(t, err)
	}

	// 验证所有实例都已注册
	discovery, err := NewDiscovery(Endpoints(clientURL))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	instances, err := discovery.GetService(ctx, "test-service")
	require.NoError(t, err)
	require.Len(t, instances, 3)
}

func TestRegistry_LeaseExpiration(t *testing.T) {
	_, clientURL := runEtcdServer(t)

	reg, err := NewRegistry(
		Endpoints(clientURL),
		TTL(2*time.Second), // 很短的 TTL
	)
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

	// 先创建 discovery（在关闭注册器之前）
	discovery, err := NewDiscovery(Endpoints(clientURL))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, discovery.Close()) })

	// 验证服务已注册
	instances, err := discovery.GetService(ctx, "test-service")
	require.NoError(t, err)
	require.Len(t, instances, 1)

	// 关闭注册器（停止续租）
	reg.Close()

	// 等待租约过期（TTL + 一些缓冲时间）
	time.Sleep(3 * time.Second)

	// 验证服务已自动过期
	instances, err = discovery.GetService(ctx, "test-service")
	require.NoError(t, err)
	require.Len(t, instances, 0, "服务应该在租约过期后自动删除")
}

func TestRegistry_Register_InvalidService(t *testing.T) {
	_, clientURL := runEtcdServer(t)

	reg, err := NewRegistry(Endpoints(clientURL))
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
}

func TestRegistry_ID(t *testing.T) {
	_, clientURL := runEtcdServer(t)

	reg, err := NewRegistry(Endpoints(clientURL))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reg.Close()) })

	require.Equal(t, RegistryID, reg.ID())
}
