// Package redis 使用 Redis 作为后端实现定位器插件
package redis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/byteweap/meta/component/locator"
	"github.com/byteweap/meta/pkg/conv"
)

// ID Redis 定位器实现标识符
const ID = "redis(hash)"

// Locator 使用 Redis 哈希结构实现 locator.Locator
type Locator struct {
	rc     redis.UniversalClient // Redis 客户端，用于哈希操作
	prefix string
}

// 确保 Locator 实现 locator.Locator 接口
var _ locator.Locator = (*Locator)(nil)

// New 使用 Redis 客户端配置创建 Redis 定位器
func New(opts redis.UniversalOptions, prefix string) *Locator {
	return newWith(redis.NewUniversalClient(&opts), prefix)
}

// newWith 使用 Redis 客户端创建 Redis 定位器
func newWith(rc redis.UniversalClient, prefix string) *Locator {
	return &Locator{
		rc:     rc,
		prefix: prefix,
	}
}

// key 生成 Redis 键
func (l *Locator) key(uid int64) string {
	if l.prefix == "" {
		return "locator:" + conv.String(uid)
	}
	return l.prefix + ":locator:" + conv.String(uid)
}

// ID 返回定位器实现标识符
func (l *Locator) ID() string {
	return ID
}

// AllNodes 返回用户所在所有服务节点
// excludes: 排除的服务
func (l *Locator) AllNodes(ctx context.Context, uid int64) (map[string]string, error) {
	return l.rc.HGetAll(ctx, l.key(uid)).Result()
}

// Node 返回用户ID当前所在的某服务某节点
func (l *Locator) Node(ctx context.Context, uid int64, service string) (string, error) {
	return l.rc.HGet(ctx, l.key(uid), service).Result()
}

// Bind 绑定用户ID到某服务某节点
func (l *Locator) Bind(ctx context.Context, uid int64, service, node string) error {
	return l.rc.HMSet(ctx, l.key(uid), service, node).Err()
}

// UnBind 如果节点匹配则解绑用户的某服务某节点
func (l *Locator) UnBind(ctx context.Context, uid int64, service, node string) error {

	current, err := l.Node(ctx, uid, service)
	if err != nil {
		return err
	}
	if current == node {
		if err = l.rc.HDel(ctx, l.key(uid), service).Err(); err != nil {
			return err
		}
	}
	return nil
}

// Close 关闭定位器
func (l *Locator) Close() error {
	return l.rc.Close()
}
