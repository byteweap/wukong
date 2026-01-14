package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/byteweap/wukong/component/registry"
	etcddiscovery "github.com/byteweap/wukong/contrib/registry/etcd"
)

var (
	mode        = flag.String("mode", "provider", "运行模式: provider(服务提供者) 或 consumer(服务消费者)")
	serviceID   = flag.String("id", "", "服务实例ID")
	serviceName = flag.String("name", "example-service", "服务名称")
	port        = flag.Int("port", 8080, "服务端口")
	etcdAddr    = flag.String("etcd", "localhost:2379", "etcd 服务地址")
)

func main() {
	flag.Parse()

	switch *mode {
	case "provider":
		runProvider()
	case "consumer":
		runConsumer()
	default:
		log.Fatalf("未知的运行模式: %s，支持的模式: provider, consumer", *mode)
	}
}

// runProvider 运行服务提供者
func runProvider() {
	// 生成服务实例ID
	instanceID := *serviceID
	if instanceID == "" {
		instanceID = fmt.Sprintf("%s-%d", *serviceName, *port)
	}

	// 创建服务注册器
	reg, err := etcddiscovery.NewRegistry(
		etcddiscovery.Endpoints(*etcdAddr),
		etcddiscovery.TTL(30*time.Second),
		etcddiscovery.Namespace("/services"),
	)
	if err != nil {
		log.Fatalf("创建注册器失败: %v", err)
	}
	defer reg.Close()

	// 创建服务实例
	service := &registry.ServiceInstance{
		ID:      instanceID,
		Name:    *serviceName,
		Version: "v1.0.0",
		Metadata: map[string]string{
			"env":    "development",
			"region": "us-east-1",
		},
		Endpoints: []string{
			fmt.Sprintf("http://127.0.0.1:%d", *port),
		},
	}

	// 注册服务
	ctx := context.Background()
	if err := reg.Register(ctx, service); err != nil {
		log.Fatalf("注册服务失败: %v", err)
	}
	log.Printf("✅ 服务已注册: %s (ID: %s, Port: %d)", *serviceName, instanceID, *port)

	// 启动 HTTP 服务器
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"service":"%s","id":"%s","status":"running","port":%d}`,
			*serviceName, instanceID, *port)
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}

	// 启动服务器
	go func() {
		log.Printf("🚀 HTTP 服务器启动在端口 %d", *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP 服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("🛑 收到停止信号，正在关闭服务...")

	// 注销服务
	if err := reg.Deregister(ctx, service); err != nil {
		log.Printf("⚠️  注销服务失败: %v", err)
	} else {
		log.Println("✅ 服务已注销")
	}

	// 关闭 HTTP 服务器
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️  HTTP 服务器关闭失败: %v", err)
	}

	log.Println("👋 服务提供者已退出")
}

// runConsumer 运行服务消费者
func runConsumer() {
	// 创建服务发现器
	registry, err := etcddiscovery.NewRegistry(
		etcddiscovery.Endpoints(*etcdAddr),
		etcddiscovery.Namespace("/services"),
	)
	if err != nil {
		log.Fatalf("创建发现器失败: %v", err)
	}
	defer registry.Close()

	ctx := context.Background()

	// 获取服务列表
	log.Printf("🔍 查找服务: %s", *serviceName)
	instances, err := registry.GetService(ctx, *serviceName)
	if err != nil {
		log.Fatalf("获取服务失败: %v", err)
	}

	if len(instances) == 0 {
		log.Printf("⚠️  未找到服务实例: %s", *serviceName)
		log.Println("💡 提示: 请先启动服务提供者 (go run main.go -mode=provider)")
		return
	}

	log.Printf("✅ 找到 %d 个服务实例:", len(instances))
	for i, instance := range instances {
		log.Printf("  [%d] ID: %s, Endpoints: %v, Metadata: %v",
			i+1, instance.ID, instance.Endpoints, instance.Metadata)
	}

	// 监听服务变更
	log.Println("\n👂 开始监听服务变更...")
	watcher, err := registry.Watch(ctx, *serviceName)
	if err != nil {
		log.Fatalf("创建监听器失败: %v", err)
	}
	defer watcher.Stop()

	// 处理服务变更
	go func() {
		for {
			instances, err := watcher.Next()
			if err != nil {
				log.Printf("⚠️  监听错误: %v", err)
				return
			}

			log.Printf("\n📢 服务变更通知: 当前有 %d 个实例", len(instances))
			for i, instance := range instances {
				log.Printf("  [%d] ID: %s, Endpoints: %v",
					i+1, instance.ID, instance.Endpoints)
			}

			// 模拟调用服务
			if len(instances) > 0 {
				endpoint := instances[0].Endpoints[0]
				log.Printf("🌐 调用服务: %s", endpoint)
				resp, err := http.Get(endpoint)
				if err != nil {
					log.Printf("❌ 调用失败: %v", err)
				} else {
					resp.Body.Close()
					log.Printf("✅ 调用成功: HTTP %s", resp.Status)
				}
			}
		}
	}()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("👋 服务消费者已退出")
}
