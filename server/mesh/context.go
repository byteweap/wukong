package mesh

import (
	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/internal/envelope"
)

type Context struct {

	// message
	req *broker.Message
	e   *envelope.Gate2MeshEnvelope

	// mesh
	m *Mesh
}

func newContext(m *Mesh, msg *broker.Message, e *envelope.Gate2MeshEnvelope) *Context {
	return &Context{
		req: msg,
		e:   e,
		m:   m,
	}
}
func (c *Context) Uid() int64 {
	if c.e == nil {
		return 0
	}
	return c.e.Uid
}

func (c *Context) Seq() uint64 {
	if c.e == nil {
		return 0
	}
	return c.e.GetMeta().GetSeq()
}

func (c *Context) App() string {
	if c.e == nil {
		return ""
	}
	return c.e.GetMeta().GetApp()
}

func (c *Context) Cmd() int32 {
	if c.e == nil {
		return 0
	}
	return c.e.GetMeta().GetCmd()
}
