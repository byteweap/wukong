# etcd 服务注册与发现示例

本示例演示如何使用 etcd 实现服务注册与发现。

## 前置要求

1. 确保本地已安装并运行 etcd 服务（默认地址: `localhost:2379`）

### 启动 etcd 服务

#### 方式一：使用 Docker Compose（推荐）

```bash
# 在 examples/registry 目录下
docker-compose up -d

# 检查 etcd 是否正常运行
docker-compose ps

# 查看日志
docker-compose logs -f etcd

# 停止 etcd
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

#### 方式二：使用 Docker 命令

```bash
docker run -d --name etcd -p 2379:2379 -p 2380:2380 \
  quay.io/coreos/etcd:v3.6.7 \
  /usr/local/bin/etcd \
  --name etcd \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://localhost:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://localhost:2380 \
  --initial-cluster etcd=http://localhost:2380

# 检查 etcd 是否运行
docker ps | grep etcd

# 使用 etcdctl 检查健康状态
docker exec etcd etcdctl endpoint health
```

#### 方式三：本地安装 etcd

如果你已经本地安装了 etcd，直接启动即可：

```bash
# 检查 etcd 是否运行
etcdctl endpoint health
```

## 运行示例

### 1. 启动服务提供者（Provider）

在第一个终端窗口运行：

```bash
# 使用默认配置
go run main.go -mode=provider

# 或指定参数
go run main.go -mode=provider -name=my-service -port=8080 -etcd=localhost:2379
```

服务提供者会：
- 注册服务到 etcd
- 启动一个 HTTP 服务器
- 提供 `/` 和 `/health` 端点
- 在收到中断信号时自动注销服务

### 2. 启动服务消费者（Consumer）

在第二个终端窗口运行：

```bash
# 使用默认配置
go run main.go -mode=consumer

# 或指定参数
go run main.go -mode=consumer -name=my-service -etcd=localhost:2379
```

服务消费者会：
- 从 etcd 发现服务实例
- 显示所有可用的服务实例
- 监听服务变更（新增/删除实例）
- 当服务变更时，自动调用服务端点

### 3. 测试服务变更

1. 启动多个服务提供者实例（不同端口）：
   ```bash
   # 终端 1
   go run main.go -mode=provider -port=8080
   
   # 终端 2
   go run main.go -mode=provider -port=8081
   
   # 终端 3
   go run main.go -mode=provider -port=8082
   ```

2. 在消费者终端中，你会看到：
   - 初始发现的所有实例
   - 当有新的提供者启动时，会收到通知
   - 当提供者停止时，也会收到通知

3. 停止一个提供者（Ctrl+C），消费者会立即检测到变更

## 命令行参数

- `-mode`: 运行模式，`provider`（服务提供者）或 `consumer`（服务消费者）
- `-id`: 服务实例ID（可选，默认自动生成）
- `-name`: 服务名称（默认: `example-service`）
- `-port`: 服务端口（默认: `8080`，仅 provider 模式）
- `-etcd`: etcd 服务地址（默认: `localhost:2379`）

## 示例输出

### 服务提供者输出
```
✅ 服务已注册: example-service (ID: example-service-8080, Port: 8080)
🚀 HTTP 服务器启动在端口 8080
🛑 收到停止信号，正在关闭服务...
✅ 服务已注销
👋 服务提供者已退出
```

### 服务消费者输出
```
🔍 查找服务: example-service
✅ 找到 3 个服务实例:
  [1] ID: example-service-8080, Endpoints: [http://127.0.0.1:8080], Metadata: map[env:development region:us-east-1]
  [2] ID: example-service-8081, Endpoints: [http://127.0.0.1:8081], Metadata: map[env:development region:us-east-1]
  [3] ID: example-service-8082, Endpoints: [http://127.0.0.1:8082], Metadata: map[env:development region:us-east-1]

👂 开始监听服务变更...

📢 服务变更通知: 当前有 2 个实例
  [1] ID: example-service-8080, Endpoints: [http://127.0.0.1:8080]
  [2] ID: example-service-8081, Endpoints: [http://127.0.0.1:8081]
🌐 调用服务: http://127.0.0.1:8080
✅ 调用成功: HTTP 200 OK
```

## 特性演示

1. **服务注册**: 服务提供者自动注册到 etcd，使用租约机制保持在线
2. **服务发现**: 服务消费者可以查询所有可用的服务实例
3. **服务监听**: 实时监听服务实例的变更（新增/删除）
4. **自动注销**: 服务提供者停止时自动注销（通过租约过期机制）
5. **优雅关闭**: 支持信号处理，优雅地关闭服务

## 注意事项

- 确保 etcd 服务正在运行
- 如果 etcd 地址不是默认的，请使用 `-etcd` 参数指定
- 服务提供者使用租约机制，如果进程异常退出，租约会在 TTL 后自动过期
- 可以同时运行多个服务提供者实例来测试负载均衡和服务发现
