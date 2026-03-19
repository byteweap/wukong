package selector

import (
	"context"

	"errors"
)

// ErrNoAvailable 表示没有可用节点
var ErrNoAvailable = errors.New("no_available_node")

// Selector 用于节点选择与负载均衡
type Selector interface {
	Rebalancer

	// Select 选择节点
	// 当 err == nil 时 selected 与 done 必须非空
	Select(ctx context.Context, opts ...SelectOption) (selected Node, done DoneFunc, err error)
}

// Rebalancer 用于节点集合更新
type Rebalancer interface {
	// Apply 在节点变更时应用全量节点
	Apply(nodes []Node)
}

// Builder 构建 Selector
type Builder interface {
	Build() Selector
}

// Node 描述服务节点
type Node interface {
	// Scheme 返回节点协议
	Scheme() string

	// Address 返回服务内唯一地址
	Address() string

	// ServiceName 返回服务名
	ServiceName() string

	// InitialWeight 返回初始权重
	// 未设置返回 nil
	InitialWeight() *int64

	// Version 返回节点版本
	Version() string

	// Metadata 返回实例元数据
	// 例如 version namespace region protocol 等
	Metadata() map[string]string
}

// DoneInfo 调用完成信息
type DoneInfo struct {
	// Err 调用错误
	Err error
	// ReplyMD 响应元数据
	ReplyMD ReplyMD

	// BytesSent 表示是否已发送请求数据
	BytesSent bool
	// BytesReceived 表示是否收到响应数据
	BytesReceived bool
}

// ReplyMD 表示响应元数据
type ReplyMD interface {
	Get(key string) string
}

// DoneFunc 调用完成回调
type DoneFunc func(ctx context.Context, di DoneInfo)
