package mesh

import (
	"sync"

	"github.com/byteweap/meta/component/broker"
)

type RPCContext struct {
	subject, reply, cmd, version string
	mesh                         *Mesh
}

var reqCtxPool = sync.Pool{
	New: func() any {
		return &RPCContext{}
	},
}

// newRPCContext 从对象池获取上下文并重置字段
func newRPCContext(mesh *Mesh, msg *broker.Message) *RPCContext {
	c := reqCtxPool.Get().(*RPCContext)
	c.reset(mesh, msg)
	return c
}

// reset 按当前消息重置上下文字段
func (ctx *RPCContext) reset(mesh *Mesh, msg *broker.Message) {
	ctx.subject = msg.Subject
	ctx.reply = msg.Reply
	ctx.cmd = msg.Header.Get("cmd")
	ctx.version = msg.Header.Get("version")
	ctx.mesh = mesh
}

// release 清理上下文字段并归还对象池
func (ctx *RPCContext) release() {
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
func (ctx *RPCContext) Subject() string {
	return ctx.subject
}

// Cmd 指令(路由)
func (ctx *RPCContext) Cmd() string {
	return ctx.cmd
}

// Version 版本
func (ctx *RPCContext) Version() string {
	return ctx.version
}
