package mesh

import "github.com/byteweap/wukong/component/broker"

type Context struct {

	// message
	req *broker.Message

	// mesh
	m *Mesh
}

func newContext(m *Mesh, msg *broker.Message) *Context {
	return &Context{
		req: msg,
		m:   m,
	}
}

func (c *Context) Subject() string {
	if c.req == nil {
		return ""
	}
	return c.req.Subject
}

func (c *Context) ReplySubject() string {
	if c.req == nil {
		return ""
	}
	return c.req.Reply
}

func (c *Context) Header() broker.Header {
	if c.req == nil {
		return nil
	}
	return c.req.Header
}

func (c *Context) Data() []byte {
	if c.req == nil {
		return nil
	}
	return c.req.Data
}
