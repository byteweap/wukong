package mesh

import (
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/envelope"
	"github.com/byteweap/wukong/internal/cluster"
	"github.com/byteweap/wukong/pkg/lang"
)

// Context 网关消息上下文
type Context struct {

	// broker message
	subject string
	reply   string // 回复的subject(由发送方传入)
	event   cluster.Event
	uid     int64

	// universal message
	seq         uint64
	fromService string
	toApp       string
	cmd         uint32
	version     uint32
	timestamp   int64

	// mesh
	mesh *Mesh
}

var ctxPool = sync.Pool{
	New: func() any {
		return &Context{}
	},
}

// newContext 从对象池获取上下文并重置字段
func newContext(mesh *Mesh, msg *broker.Message, e *envelope.IMessage) *Context {
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
	c.event = ""
	c.uid = 0

	c.seq = 0
	c.fromService = ""
	c.toApp = ""
	c.cmd = 0
	c.version = 0
	c.mesh = nil
	ctxPool.Put(c)
}

// reset 按当前消息重置上下文字段
func (c *Context) reset(mesh *Mesh, msg *broker.Message, e *envelope.IMessage) {

	c.mesh = mesh

	// broker
	c.subject = msg.Subject
	c.reply = cluster.GetReplyBy(msg.Header)
	c.fromService = cluster.GetFromServiceBy(msg.Header)
	c.toApp = cluster.GetToServiceBy(msg.Header)
	c.event = cluster.GetEventBy(msg.Header)
	c.uid = cluster.GetUidBy(msg.Header)

	if e == nil || e.GetHeader() == nil {
		c.seq = 0
		c.cmd = 0
		c.version = 0
		c.timestamp = 0
		return
	}
	header := e.GetHeader()
	c.seq = header.GetSeq()
	c.cmd = header.GetCmd()
	c.version = header.GetVersion()
	c.timestamp = header.GetTimestamp()
}

// Uid 返回用户 ID
func (c *Context) Uid() int64 {
	return c.uid
}

// Seq 返回消息序列号
func (c *Context) Seq() uint64 {
	return c.seq
}

// FromService 返回应用标识
func (c *Context) FromService() string {
	return c.fromService
}

// ToService 返回应用标识
func (c *Context) ToService() string {
	return c.toApp
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

// Timestamp 返回消息时间戳
func (c *Context) Timestamp() int64 {
	return c.timestamp
}

// Copy 复制
// 当在新的goroutine中使用context时,应使用Copy方法复制context
func (c *Context) Copy() *Context {
	if c == nil {
		return &Context{}
	}
	return &Context{
		subject:     c.subject,
		reply:       c.reply,
		event:       c.event,
		uid:         c.uid,
		seq:         c.seq,
		fromService: c.fromService,
		toApp:       c.toApp,
		cmd:         c.cmd,
		version:     c.version,
		timestamp:   c.timestamp,
		mesh:        c.mesh,
	}
}

// OkResp 返回成功响应
func (c *Context) OkResp(args ...proto.Message) {

	out := &envelope.OMessage{
		Header: &envelope.Header{
			Seq:       c.Seq(),
			Cmd:       c.Cmd(),
			Version:   c.Version(),
			Timestamp: time.Now().UnixMilli(),
		},
		Service: c.mesh.appName,
		MsgType: envelope.MsgType_RESPONSE,
	}

	var err error
	if len(args) > 0 {
		out.Payload, err = proto.Marshal(args[0])
		if err != nil {
			log.Errorf("[mesh].[OkResponse] marshal payload error, err: %v", err)
			return
		}
	}
	bytes, err := proto.Marshal(out)
	if err != nil {
		log.Errorf("[mesh].[OkResponse] marshal message error, err: %v", err)
		return
	}
	// 发送
	if c.reply == "" {
		log.Errorf("[mesh].[OkResponse] reply subject is empty")
		return
	}
	if err = c.mesh.sendMessage(c.reply, c.fromService, bytes, c.uid); err != nil {
		log.Errorf("[mesh].[OkResponse] send message error, subject: %v, err: %v", c.reply, err)
		return
	}
}

// ErrResp 返回错误响应
func (c *Context) ErrResp(code int, args ...string) {

	tip := lang.If(len(args) > 0, args[0], "mesh internal error")

	out := &envelope.OMessage{
		Header: &envelope.Header{
			Seq:       c.Seq(),
			Cmd:       c.Cmd(),
			Version:   c.Version(),
			Timestamp: time.Now().UnixMilli(),
		},
		Service: c.mesh.appName,
		MsgType: envelope.MsgType_RESPONSE,
		Result: &envelope.Code{
			Code: int32(code),
			Tip:  tip,
		},
	}
	bytes, err := proto.Marshal(out)
	if err != nil {
		log.Errorf("[mesh].[ErrResponse] marshal message error, err: %v", err)
		return
	}
	// 发送
	if c.reply == "" {
		log.Errorf("[mesh].[ErrResponse] reply subject is empty")
		return
	}
	if err = c.mesh.sendMessage(c.reply, c.fromService, bytes, c.uid); err != nil {
		log.Errorf("[mesh].[ErrResponse] send message error, subject: %v, err: %v", c.reply, err)
	}
}

func (c *Context) Broadcast() {
	// todo
}
