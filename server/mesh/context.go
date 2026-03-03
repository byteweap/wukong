package mesh

import (
	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/internal/envelope"
)

// Context 网关消息上下文
type Context struct {

	// broker message
	subject, reply string
	header         broker.Header

	// universal message
	seq uint64
	app string
	cmd int32
	uid int64

	// mesh
	mesh *Mesh
}

func newContext(mesh *Mesh, msg *broker.Message, e *envelope.Gate2MeshEnvelope) *Context {
	return &Context{
		subject: msg.Subject,
		reply:   msg.Reply,
		header:  msg.Header,
		seq:     e.GetMeta().GetSeq(),
		app:     e.GetMeta().GetApp(),
		cmd:     e.GetMeta().GetCmd(),
		uid:     e.GetUid(),
		mesh:    mesh,
	}
}

func (c *Context) Uid() int64 {
	return c.uid
}

func (c *Context) Seq() uint64 {
	return c.seq
}

func (c *Context) App() string {
	return c.app
}

func (c *Context) Cmd() int32 {
	return c.cmd
}

func (c *Context) Subject() string {
	return c.subject
}

func (c *Context) ReplySubject() string {
	return c.reply
}

func (c *Context) Header() broker.Header {
	return c.header
}
