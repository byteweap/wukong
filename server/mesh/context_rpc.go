package mesh

import (
	"sync"

	"github.com/byteweap/wukong/component/broker"
)

type RpcContext struct {
	subject, reply, cmd, version string
	mesh                         *Mesh
}

var reqCtxPool = sync.Pool{
	New: func() any {
		return &RpcContext{}
	},
}

// newRpcContext 从对象池获取上下文并重置字段
func newRpcContext(mesh *Mesh, msg *broker.Message) *RpcContext {
	c := reqCtxPool.Get().(*RpcContext)
	c.reset(mesh, msg)
	return c
}

// reset 按当前消息重置上下文字段
func (ctx *RpcContext) reset(mesh *Mesh, msg *broker.Message) {
	ctx.subject = msg.Subject
	ctx.reply = msg.Reply
	ctx.cmd = msg.Header.Get("cmd")
	ctx.version = msg.Header.Get("version")
	ctx.mesh = mesh
}

// release 清理上下文字段并归还对象池
func (ctx *RpcContext) release() {
	if ctx == nil {
		return
	}
	ctx.subject = ""
	ctx.reply = ""
	ctx.cmd = ""
	ctx.version = ""
	ctx.mesh = nil
	reqCtxPool.Put(ctx)
}

// Subject 获取当前请求的主题
func (ctx *RpcContext) Subject() string {
	return ctx.subject
}

// Cmd 指令(路由)
func (ctx *RpcContext) Cmd() string {
	return ctx.cmd
}

// Version 版本
func (ctx *RpcContext) Version() string {
	return ctx.version
}
