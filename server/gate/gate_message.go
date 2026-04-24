package gate

import (
	"errors"
	"net/http"

	"github.com/byteweap/meta/component/broker"
	"github.com/byteweap/meta/component/log"
	"github.com/byteweap/meta/encoding/proto"
	"github.com/byteweap/meta/envelope"
	"github.com/byteweap/meta/internal/cluster"
	"github.com/byteweap/meta/pkg/conv"
)

// handlerRequestReplyMessage 来自其它服务的(request-reply)消息
func (g *Gate) handleRequestReplyMessage(msg *broker.Message) {
	if msg == nil {
		return
	}
	if err := g.replyError(msg, http.StatusNotImplemented, "gate request-reply is not implemented"); err != nil {
		log.Errorf("[websocket] request-reply unsupported, reply error: %v", err)
	}
}

// handlerPubSubMessage 来自Mesh服务的(pub-sub)消息
func (g *Gate) handlePubSubMessage(msg *broker.Message) {
	// 2. 直接回复给玩家的消息
	uid := cluster.GetUidBy(msg.Header)
	if uid <= 0 {
		log.Errorf("[websocket] reply2player get uid error, uid: %v", uid)
		return
	}
	session, ok := g.sessions.get(uid)
	if !ok {
		log.Errorf("[websocket] reply2player get session error, uid: %v", uid)
		return
	}
	if err := session.WriteBinary(msg.Data); err != nil {
		log.Errorf("[websocket] reply2player write binary error, uid: %v, err: %v", uid, err)
		return
	}
	log.Debugf("[websocket] reply2player success, uid: %v", uid)
}

// 处理来自其它服务的消息
func (g *Gate) handleMessage(msg *broker.Message) {
	if msg.Reply != "" {
		g.handleRequestReplyMessage(msg)
	} else {
		g.handlePubSubMessage(msg)
	}
}

func (g *Gate) replyError(reqMsg *broker.Message, code int, tip string) error {
	if reqMsg == nil {
		return errors.New("request message is nil")
	}
	if reqMsg.Reply == "" {
		return errors.New("reply subject is empty")
	}
	header := reqMsg.Header
	if header == nil {
		header = broker.Header{}
	}
	header.Set("code", conv.String(code))
	header.Set("tip", tip)
	return g.opts.broker.Reply(g.ctx, reqMsg, nil, broker.ReplyHeader(header))
}

// 业务消息分发至 mesh
func (g *Gate) dispatch(uid int64, e *envelope.IMessage) {

	if e == nil {
		log.Errorf("[websocket] dispatch error, envelope is nil")
		return
	}
	var (
		toService = e.GetService()
		loc, bro  = g.opts.locator, g.opts.broker
	)

	curNodeID, err := loc.Node(g.ctx, uid, toService)
	if err != nil {
		log.Errorf("[websocket] dispatch | get mesh node error, uid: %v, toService: %v, err: %v", uid, toService, err)
		return
	}

	data, err := proto.Marshal(e)
	if err != nil {
		log.Errorf("[websocket] dispatch | marshal to mesh data error: %v", err)
		return
	}
	nodeID := curNodeID
	if curNodeID == "" {
		sel, err := g.ensure(toService)
		if err != nil {
			log.Errorf("[websocket] dispatch | get mesh node error, uid: %v, toService: %v, err: %v", uid, toService, err)
			return
		}
		node, err := sel.Select("")
		if err != nil {
			log.Errorf("[websocket] dispatch | select mesh node error, uid: %v, toService: %v, err: %v", uid, toService, err)
			return
		}
		nodeID = node.ID()
	}
	// 构建消息头
	var (
		reply  = g.Subject(toService) // 回复主题
		header = cluster.BuildHeader(uid, cluster.Event_Business, reply, g.appName, toService)
	)
	// 发布消息到 Mesh
	subject := cluster.Subject(g.opts.prefix, g.appName, toService, nodeID)
	if err = bro.Pub(g.ctx, subject, data, broker.PubHeader(header)); err != nil {
		log.Errorf("[websocket] dispatch error, uid: %v, subject: %v, err: %v", uid, subject, err)
		return
	}
	log.Debugf("[websocket] dispatch success, uid: %v, subject: %v", uid, subject)
}

// 广播系统事件
func (g *Gate) broadcastEvent(uid int64, event cluster.Event) {

	// 获取玩家当前所在所有节点
	snMap, err := g.opts.locator.AllNodes(g.ctx, uid)
	if err != nil {
		log.Errorf("[websocket] broadcast event, get all nodes error, uid: %v, event: %v, err: %v", uid, event, err)
		return
	}
	for service, node := range snMap {
		if service == g.appName { // 不包括 gate
			continue
		}
		// 发布消息到 Mesh
		var (
			header  = cluster.BuildHeader(uid, event, g.Subject(service), g.appName, service)
			subject = cluster.Subject(g.opts.prefix, g.appName, service, node)
		)
		if err = g.opts.broker.Pub(g.ctx, subject, nil, broker.PubHeader(header)); err != nil {
			log.Errorf("[websocket] broadcast event error, uid: %v, subject: %v, err: %v", uid, subject, err)
			return
		}
		log.Debugf("[websocket] broadcast event success, uid: %v, subject: %v, event: %v", uid, subject, event)
	}
}
