// Package redis 使用 Redis 作为后端实现定位器插件
package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/byteweap/wukong/component/locator"
)

// ID Redis 定位器实现标识符
const ID = "redis(hash)"

// Locator 使用 Redis 哈希结构实现 locator.Locator
type Locator struct {
	rc                redis.UniversalClient // Redis 客户端，用于哈希操作
	keyFormat         string                // Redis 键格式化字符串
	gateNodeFieldName string                // 哈希结构中网关节点字段名
	gameNodeFieldName string                // 哈希结构中游戏节点字段名
}

// 确保 Locator 实现 locator.Locator 接口
var _ locator.Locator = (*Locator)(nil)

// New 使用 Redis 客户端配置创建 Redis 定位器
func New(opts redis.UniversalOptions, keyFormat, gateNodeFieldName, gameNodeFieldName string) *Locator {
	return newWith(redis.NewUniversalClient(&opts), keyFormat, gateNodeFieldName, gameNodeFieldName)
}

// newWith 使用 Redis 客户端创建 Redis 定位器
func newWith(rc redis.UniversalClient, keyFormat, gateNodeFieldName, gameNodeFieldName string) *Locator {
	return &Locator{
		rc:                rc,
		keyFormat:         keyFormat,
		gateNodeFieldName: gateNodeFieldName,
		gameNodeFieldName: gameNodeFieldName,
	}
}

// ID 返回定位器实现标识符
func (l *Locator) ID() string {
	return ID
}

// Gate 返回用户ID对应的网关节点
func (l *Locator) Gate(ctx context.Context, uid int64) (string, error) {

	key := fmt.Sprintf(l.keyFormat, uid)
	return l.rc.HGet(ctx, key, l.gateNodeFieldName).Result()
}

// BindGate 绑定用户ID到网关节点
func (l *Locator) BindGate(ctx context.Context, uid int64, node string) error {

	key := fmt.Sprintf(l.keyFormat, uid)
	return l.rc.HMSet(ctx, key, l.gateNodeFieldName, node).Err()
}

// UnBindGate 如果节点匹配则解绑用户ID的网关节点
func (l *Locator) UnBindGate(ctx context.Context, uid int64, node string) error {

	current, err := l.Gate(ctx, uid)
	if err != nil {
		return err
	}

	if current == node {
		key := fmt.Sprintf(l.keyFormat, uid)
		if err = l.rc.HMSet(ctx, key, l.gateNodeFieldName, "").Err(); err != nil {
			return err
		}
	}
	return nil
}

// Game 返回用户ID对应的游戏节点
func (l *Locator) Game(ctx context.Context, uid int64) (string, error) {

	key := fmt.Sprintf(l.keyFormat, uid)
	return l.rc.HGet(ctx, key, l.gameNodeFieldName).Result()
}

// BindGame 绑定用户ID到游戏节点
func (l *Locator) BindGame(ctx context.Context, uid int64, node string) error {

	key := fmt.Sprintf(l.keyFormat, uid)
	return l.rc.HMSet(ctx, key, l.gameNodeFieldName, node).Err()
}

// UnBindGame 如果节点匹配则解绑用户ID的游戏节点
func (l *Locator) UnBindGame(ctx context.Context, uid int64, node string) error {

	current, err := l.Game(ctx, uid)
	if err != nil {
		return err
	}

	if current == node {
		key := fmt.Sprintf(l.keyFormat, uid)
		if err = l.rc.HMSet(ctx, key, l.gameNodeFieldName, "").Err(); err != nil {
			return err
		}
	}
	return nil
}

// Close 关闭定位器
func (l *Locator) Close() error {
	return l.rc.Close()
}
