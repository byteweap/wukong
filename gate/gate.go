package gate

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/component/network"
	"github.com/byteweap/wukong/component/registry"
	"github.com/byteweap/wukong/pkg/xnet"
)

// Gate websocket 网关
type Gate struct {
	opts           *options // 配置选项
	ctx            context.Context
	cancel         context.CancelFunc
	sessionManager *SessionManager // 会话管理器
	mu             sync.Mutex
	instance       *registry.ServiceInstance // 服务实例
}

// New 创建新的网关服务器实例
func New(opts ...Option) *Gate {

	// 应用配置选项
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	if o.logger != nil {
		log.SetLogger(o.logger)
	}

	ctx, cancel := context.WithCancel(o.ctx)
	return &Gate{
		ctx:            ctx,
		cancel:         cancel,
		opts:           o,
		sessionManager: NewSessionManager(),
	}
}

// Run 启动网关服务器
func (g *Gate) Run() error {

	// 验证配置选项
	if err := g.validate(); err != nil {
		return err
	}

	// 构建服务实例
	if err := g.buildInstance(); err != nil {
		return err
	}

	// 初始化网络配置
	g.setupNetwork()

	// 注册服务
	if err := g.registerService(); err != nil {
		return err
	}

	// 启动网络服务器 (阻塞)
	g.opts.netServer.Start()

	// 停止网关服务器
	err := g.Stop()
	if err != nil {
		return err
	}

	return nil
}

// Stop 关闭网关服务器
func (g *Gate) Stop() error {

	if g.opts.netServer != nil {
		g.opts.netServer.Stop()
	}
	if g.opts.broker != nil {
		if err := g.opts.broker.Close(); err != nil {
			return fmt.Errorf("broker close error: %w", err)
		}
	}
	if g.opts.locator != nil {
		if err := g.opts.locator.Close(); err != nil {
			return fmt.Errorf("locator close error: %w", err)
		}
	}

	if err := g.sessionManager.Close(); err != nil {
		return fmt.Errorf("session manager close error: %w", err)
	}

	if err := g.unregisterService(); err != nil {
		return fmt.Errorf("unregister service error: %w", err)
	}
	log.Info("gate stop successfully")
	return nil
}

// setupNetwork 初始化网络服务器配置
func (g *Gate) setupNetwork() {

	// todo
}

// handlerBinaryMessage 处理接收到的二进制消息
func (g *Gate) handlerBinaryMessage(_ network.Conn, msg []byte) {
	fmt.Println("Gate receive binary message: ", msg)
	// todo
}

// buildInstance 构建服务实例
func (g *Gate) buildInstance() error {

	app := g.opts.app

	endpoints := make([]string, 0, len(app.endpoints))
	for _, e := range app.endpoints {
		endpoints = append(endpoints, e.String())
	}
	if len(endpoints) == 0 {
		ip, err := xnet.ExternalIP()
		if err != nil {
			return fmt.Errorf("get external ip error: %w", err)
		}
		e := fmt.Sprintf("ws://%s%s", ip, g.opts.netServer.Addr())
		endpoints = append(endpoints, e)
	}
	instance := &registry.ServiceInstance{
		ID:        app.id,
		Name:      app.name,
		Version:   app.version,
		Metadata:  app.metadata,
		Endpoints: endpoints,
	}
	g.mu.Lock()
	g.instance = instance
	g.mu.Unlock()

	return nil
}

func (g *Gate) validate() error {
	o := g.opts

	if o.netServer == nil {
		return errors.New("network server is not set")
	}
	if o.locator == nil {
		return errors.New("locator is not set")
	}
	if o.broker == nil {
		return errors.New("broker is not set")
	}
	if o.registry == nil {
		return errors.New("registry is not set")
	}

	return nil
}

// registerService 注册服务
func (g *Gate) registerService() error {

	if g.opts.registry == nil {
		return nil // 如果注册器为空，则不进行注册
	}

	g.mu.Lock()
	instance := g.instance
	g.mu.Unlock()

	if instance == nil {
		return fmt.Errorf("service instance is nil")
	}

	ctx, cancel := context.WithTimeout(g.ctx, g.opts.registryTimeout)
	defer cancel()

	return g.opts.registry.Register(ctx, instance)
}

// unregisterService 注销服务
func (g *Gate) unregisterService() error {

	if g.opts.registry == nil {
		return nil // 如果注册器为空，则不进行注销
	}

	g.mu.Lock()
	instance := g.instance
	g.mu.Unlock()
	if instance == nil {
		return fmt.Errorf("service instance is nil")
	}

	ctx, cancel := context.WithTimeout(g.ctx, g.opts.registryTimeout)
	defer cancel()
	return g.opts.registry.Deregister(ctx, instance)
}
