package selector

import (
	"context"
)

type peerKey struct{}

// Peer 表示一次 RPC 的对端信息
// 例如地址与认证信息
type Peer struct {
	// Node 是对端节点
	Node Node
}

// NewPeerContext 创建带 Peer 的上下文
func NewPeerContext(ctx context.Context, p *Peer) context.Context {
	return context.WithValue(ctx, peerKey{}, p)
}

// FromPeerContext 从上下文获取 Peer
func FromPeerContext(ctx context.Context) (p *Peer, ok bool) {
	p, ok = ctx.Value(peerKey{}).(*Peer)
	return
}
