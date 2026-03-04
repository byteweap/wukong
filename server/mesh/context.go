package mesh

import (
	"maps"
	"sync"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/internal/envelope"
)

// Context 网关消息上下文
type Context struct {

	// broker message
	subject string
	header  broker.Header

	// universal message
	seq     uint64
	app     string
	cmd     uint32
	version uint32
	uid     int64

	// mesh
	mesh *Mesh
}

var ctxPool = sync.Pool{
	New: func() any {
		return &Context{}
	},
}

// newContext 从对象池获取上下文并重置字段
func newContext(mesh *Mesh, msg *broker.Message, e *envelope.Gate2MeshEnvelope) *Context {
	c := ctxPool.Get().(*Context)
	c.reset(mesh, msg, e)
	return c
}

// release 清理上下文字段并归还对象池
func (c *Context) release() {
	if c == nil {
		return
	}
	c.subject = ""
	c.header = nil
	c.seq = 0
	c.app = ""
	c.cmd = 0
	c.version = 0
	c.uid = 0
	c.mesh = nil
	ctxPool.Put(c)
}

// reset 按当前消息重置上下文字段
func (c *Context) reset(mesh *Mesh, msg *broker.Message, e *envelope.Gate2MeshEnvelope) {
	c.subject = msg.Subject
	c.header = msg.Header
	c.mesh = mesh
	if e == nil {
		c.seq = 0
		c.app = ""
		c.cmd = 0
		c.version = 0
		c.uid = 0
		return
	}
	c.uid = e.GetUid()
	meta := e.GetMeta()
	if meta == nil {
		c.seq = 0
		c.app = ""
		c.cmd = 0
		c.version = 0
		return
	}
	c.seq = meta.GetSeq()
	c.app = meta.GetApp()
	c.cmd = meta.GetCmd()
	c.version = meta.GetVersion()
}

// Uid 返回用户 ID
func (c *Context) Uid() int64 {
	return c.uid
}

// Seq 返回消息序列号
func (c *Context) Seq() uint64 {
	return c.seq
}

// App 返回应用标识
func (c *Context) App() string {
	return c.app
}

// Cmd 返回路由命令字
func (c *Context) Cmd() uint32 {
	return c.cmd
}

// Version 返回消息版本号
func (c *Context) Version() uint32 {
	return c.version
}

// Subject 返回 broker 主题
func (c *Context) Subject() string {
	return c.subject
}

// Header 返回 broker 消息头
func (c *Context) Header() broker.Header {
	return c.header
}

// Copy 复制
// 当在新的goroutine中使用context时,应使用Copy方法复制context
func (c *Context) Copy() *Context {
	if c == nil {
		return &Context{}
	}
	return &Context{
		subject: c.subject,
		header:  copyHeader(c.header),
		seq:     c.seq,
		app:     c.app,
		cmd:     c.cmd,
		uid:     c.uid,
		mesh:    c.mesh,
	}
}

// copyHeader 深拷贝 header，避免异步共享底层引用
func copyHeader(h broker.Header) broker.Header {
	if h == nil {
		return nil
	}
	cp := make(broker.Header, len(h))
	maps.Copy(cp, h)
	return cp
}
