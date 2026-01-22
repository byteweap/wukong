package gate

import (
	"context"
	"fmt"
	"sync"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/locator"
	"github.com/byteweap/wukong/component/network"
	"github.com/byteweap/wukong/component/registry"
	"github.com/google/uuid"
)

// Gate websocket 网关
type Gate struct {
	ctx    context.Context
	cancel context.CancelFunc

	opts           *options          // 配置选项
	netServer      network.Server    // 网络服务器（WebSocket）
	locator        locator.Locator   // 玩家位置定位器
	broker         broker.Broker     // 消息传输代理
	sessionManager *SessionManager   // 会话管理器
	registry       registry.Registry // 服务注册与发现器

	mu       sync.Mutex
	instance *registry.ServiceInstance // 服务实例
}

// New 创建新的网关服务器实例
func New(opts ...Option) (*Gate, error) {

	// 应用配置选项
	o := defaultOptions()

	if id, err := uuid.NewUUID(); err == nil {
		o.application.ID = id.String()
	}
	for _, opt := range opts {
		opt(o)
	}

	return &Gate{
		opts:           o,
		sessionManager: NewSessionManager(),
	}, nil
}

// Run 启动网关服务器
func (g *Gate) Run() error {

	// 构建服务实例
	g.buildInstance()

	// 初始化网络配置
	g.setupNetwork()

	// 注册服务
	if err := g.registerService(); err != nil {
		return err
	}

	// 启动网络服务器 (阻塞)
	g.netServer.Start()

	// 停止网关服务器
	err := g.Stop()
	if err != nil {
		return err
	}

	return nil
}

// Stop 关闭网关服务器
func (g *Gate) Stop() error {

	g.netServer.Stop()

	if err := g.broker.Close(); err != nil {
		return fmt.Errorf("broker close error: %w", err)
	}

	if err := g.locator.Close(); err != nil {
		return fmt.Errorf("locator close error: %w", err)
	}

	if err := g.sessionManager.Close(); err != nil {
		return fmt.Errorf("session manager close error: %w", err)
	}

	return nil
}

// setupNetwork 初始化网络服务器配置
func (g *Gate) setupNetwork() {

}

// handlerBinaryMessage 处理接收到的二进制消息
func (g *Gate) handlerBinaryMessage(_ network.Conn, msg []byte) {
	fmt.Println("Gate receive binary message: ", msg)
}

// buildInstance 构建服务实例
func (g *Gate) buildInstance() {

	app := g.opts.application
	instance := &registry.ServiceInstance{
		ID:        app.ID,
		Name:      app.Name,
		Version:   app.Version,
		Metadata:  app.Metadata,
		Endpoints: []string{app.Addr},
	}

	g.mu.Lock()
	g.instance = instance
	g.mu.Unlock()
}

// registerService 注册服务
func (g *Gate) registerService() error {

	g.mu.Lock()
	instance := g.instance
	g.mu.Unlock()
	if g.registry == nil {
		// 如果注册器为空，则不进行注册
		return nil
	}
	if instance == nil {
		return fmt.Errorf("service instance is nil")
	}

	ctx, cancel := context.WithTimeout(g.ctx, g.opts.registryTimeout)
	defer cancel()

	return g.registry.Register(ctx, instance)
}

// unregisterService 注销服务
func (g *Gate) unregisterService() error {

	g.mu.Lock()
	instance := g.instance
	g.mu.Unlock()

	if g.registry == nil {
		// 如果注册器为空，则不进行注销
		return nil
	}
	if instance == nil {
		return fmt.Errorf("service instance is nil")
	}

	ctx, cancel := context.WithTimeout(g.ctx, g.opts.registryTimeout)
	defer cancel()

	return g.registry.Deregister(ctx, instance)
}
