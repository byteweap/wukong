# Nacos 服务注册与发现示例

本示例演示如何使用 Nacos 实现服务注册与发现。

## 前置要求

1. 确保本地已安装并运行 Nacos 服务（默认地址: `127.0.0.1:8848`）

## 运行示例

### 1. 启动服务提供者（Provider）

在第一个终端窗口运行：

```bash
# 使用默认配置
go run main.go -mode=provider

# 或指定参数
go run main.go -mode=provider -name=my-service -port=8080 -nacos=127.0.0.1:8848
```

服务提供者会：
- 注册服务到 Nacos
- 启动一个 HTTP 服务器
- 提供 `/` 和 `/health` 端点
- 在收到中断信号时自动注销服务
- 使用心跳机制保持服务在线（Ephemeral 实例）

### 2. 启动服务消费者（Consumer）

在第二个终端窗口运行：

```bash
# 使用默认配置
go run main.go -mode=consumer

# 或指定参数
go run main.go -mode=consumer -name=my-service -nacos=127.0.0.1:8848
```

服务消费者会：
- 从 Nacos 发现服务实例
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

### 4. 使用不同的命名空间和分组

```bash
# 在命名空间 "dev" 中注册服务
go run main.go -mode=provider -namespace=dev -group=web-services

# 在命名空间 "dev" 中发现服务
go run main.go -mode=consumer -namespace=dev -group=web-services
```

## 命令行参数

- `-mode`: 运行模式，`provider`（服务提供者）或 `consumer`（服务消费者）
- `-id`: 服务实例ID（可选，默认自动生成）
- `-name`: 服务名称（默认: `example-service`）
- `-port`: 服务端口（默认: `8080`，仅 provider 模式）
- `-nacos`: Nacos 服务地址（默认: `127.0.0.1:8848`）
- `-namespace`: Nacos 命名空间（默认: `public`）
- `-group`: 服务分组（默认: `DEFAULT_GROUP`）

## 示例输出

### 服务提供者输出
```
✅ 服务已注册: example-service (ID: example-service-8080, Port: 8080, Namespace: public, Group: DEFAULT_GROUP)
🚀 HTTP 服务器启动在端口 8080
🛑 收到停止信号，正在关闭服务...
✅ 服务已注销
👋 服务提供者已退出
```

### 服务消费者输出
```
🔍 查找服务: example-service (Namespace: public, Group: DEFAULT_GROUP)
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

1. **服务注册**: 服务提供者自动注册到 Nacos，使用临时实例（Ephemeral）和心跳机制保持在线
2. **服务发现**: 服务消费者可以查询所有可用的服务实例
3. **服务监听**: 实时监听服务实例的变更（新增/删除）
4. **自动注销**: 服务提供者停止时自动注销（通过心跳超时机制）
5. **优雅关闭**: 支持信号处理，优雅地关闭服务
6. **命名空间隔离**: 支持使用不同的命名空间隔离服务
7. **服务分组**: 支持使用服务分组组织服务

## Nacos 控制台

访问 Nacos 控制台可以查看和管理服务：

1. 打开浏览器访问: http://localhost:8848/nacos
2. 使用默认账号登录: `nacos` / `nacos`
3. 在"服务管理" -> "服务列表"中查看注册的服务
4. 可以查看服务实例的详细信息、健康状态等

## 注意事项

- 确保 Nacos 服务正在运行
- 如果 Nacos 地址不是默认的，请使用 `-nacos` 参数指定
- 服务提供者使用临时实例（Ephemeral=true），如果进程异常退出，实例会在心跳超时后自动下线
- 可以同时运行多个服务提供者实例来测试负载均衡和服务发现
- 使用不同的命名空间和分组可以实现服务隔离和组织
- 心跳间隔默认 5 秒，可以通过代码中的 `BeatInterval` 选项调整

## 故障排查

### Nacos 启动失败

如果遇到 Nacos 启动失败，常见原因和解决方案：

1. **数据库连接错误**：如果看到 MySQL 连接错误
   - 确保 MySQL 服务已启动并健康：`docker-compose ps`
   - 检查 MySQL 日志：`docker-compose logs mysql`
   - 确认数据库初始化完成：等待 MySQL 健康检查通过
   - 检查环境变量是否正确：`MYSQL_SERVICE_HOST=mysql`（使用服务名）
   - 如果使用外部 MySQL，确保网络可达

2. **MySQL 初始化失败**
   - 检查 `init.sql` 文件是否存在
   - 查看 MySQL 启动日志确认初始化是否成功
   - 如果初始化失败，可以手动执行 SQL：
     ```bash
     docker exec -i nacos-mysql mysql -uroot -pnacos nacos < init.sql
     ```

2. **端口冲突**：如果 8848 端口已被占用
   ```bash
   # 检查端口占用
   netstat -ano | findstr :8848  # Windows
   lsof -i :8848                  # Linux/Mac
   
   # 修改 docker-compose.yml 中的端口映射
   ports:
     - "8849:8848"  # 使用其他端口
   ```

3. **权限问题**：确保 Docker 有足够权限访问卷
   ```bash
   # 清理旧的卷和数据
   docker-compose down -v
   docker volume prune
   ```

4. **内存不足**：如果容器内存不足，可以调整 JVM 参数
   ```yaml
   environment:
     - JVM_XMS=256m
     - JVM_XMX=256m
     - JVM_XMN=128m
   ```

### 查看服务日志

```bash
# 查看 Nacos 容器日志
docker-compose logs -f nacos

# 查看 MySQL 容器日志
docker-compose logs -f mysql

# 查看所有服务日志
docker-compose logs -f

# 查看 Nacos 详细日志
docker exec nacos-server tail -f /home/nacos/logs/nacos.log

# 进入 MySQL 容器检查数据库
docker exec -it nacos-mysql mysql -uroot -pnacos
# 然后执行: USE nacos; SHOW TABLES;
```

### 重置所有服务

如果遇到数据问题，可以重置所有服务：

```bash
# 停止并删除容器和数据卷（会清除所有数据）
docker-compose down -v

# 重新启动（会重新初始化数据库）
docker-compose up -d
```

### 仅重置 Nacos（保留 MySQL 数据）

如果只想重置 Nacos 而保留 MySQL 数据：

```bash
# 停止并删除 Nacos 容器
docker-compose stop nacos
docker-compose rm -f nacos

# 重新启动 Nacos
docker-compose up -d nacos
```

### 检查服务状态

```bash
# 检查所有服务状态
docker-compose ps

# 检查服务健康状态
docker-compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"

# 测试 Nacos API
curl http://localhost:8848/nacos/v1/console/health/liveness

# 测试 MySQL 连接
docker exec nacos-mysql mysqladmin ping -h localhost -uroot -pnacos
```

## 与 etcd 实现的区别

| 特性 | etcd | Nacos |
|------|------|-------|
| 心跳机制 | 租约续期 | 主动发送心跳 |
| 实例类型 | 永久实例 | 临时实例（Ephemeral） |
| 命名空间 | Key 前缀 | 独立命名空间 |
| 服务分组 | 不支持 | 支持 |
| 控制台 | 无 | 提供 Web 控制台 |
| 健康检查 | 租约过期 | 心跳超时 |
