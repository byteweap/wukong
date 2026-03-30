package gate

import (
	"github.com/olahol/melody"

	"github.com/byteweap/meta/component/log"
	"github.com/byteweap/meta/encoding/proto"
	"github.com/byteweap/meta/envelope"
	"github.com/byteweap/meta/internal/cluster"
)

// 连接建立时调用
func (g *Gate) handleConnect(s *melody.Session) {

	req, addr := s.Request, s.RemoteAddr()

	uid := g.opts.userIdExtractor(req)
	if uid <= 0 {
		_ = s.Write([]byte("uid is required"))
		_ = s.Close()
		return
	}
	// 注册会话
	session, ok := g.sessions.get(uid)
	if ok {
		log.Warnf("[websocket] connection exists: uid: %v, close old connection", uid)
		_ = session.Close()
	}
	s.Set("uid", uid)
	g.sessions.register(uid, s)

	log.Infof("[websocket] new connection success, uid: %v, %s", uid, addr.String())

	loc := g.opts.locator

	// 绑定网关
	if err := loc.Bind(g.ctx, uid, g.appName, g.appID); err != nil {
		log.Errorf("[websocket] new connection success, bind gate error, uid: %v, err: %v", uid, err)
		return
	}

	// 广播 上线、重连 事件到上游服务
	event := cluster.Event_Online
	if ok {
		event = cluster.Event_Reconnect
	}
	g.broadcastEvent(uid, event)
}

// 连接断开时调用
func (g *Gate) handleDisconnect(s *melody.Session) {

	uids, ok := s.Get("uid")
	if !ok {
		log.Error("[websocket] connection disconnect error, session not contains uid key")
		return
	}
	uid := uids.(int64)

	// 注销会话
	curSession, yes := g.sessions.get(uid)
	if !yes {
		log.Errorf("[websocket] connection disconnect error, uid: %v not found", uid)
		return
	}
	if curSession != s {
		log.Warnf("[websocket] connection disconnect error, uid: %v session not match", uid)
		return
	}
	g.sessions.unregister(uid)

	log.Infof("[websocket] connection disconnect success, uid: %v", uid)

	// 解绑网关
	if err := g.opts.locator.UnBind(g.ctx, uid, g.appName, g.appID); err != nil {
		log.Errorf("[websocket] connection disconnect success, unbind gate error, uid: %v, err: %v", uid, err)
	}

	// 广播掉线事件到上游服务
	g.broadcastEvent(uid, cluster.Event_Offline)
}

// 接收到文本消息时调用
func (g *Gate) handleTextMessage(s *melody.Session, msg []byte) {
	// todo
}

// 接收到二进制消息时调用
func (g *Gate) handleBinaryMessage(s *melody.Session, msg []byte) {

	meta := &envelope.IMessage{}
	if err := proto.Unmarshal(msg, meta); err != nil {
		log.Errorf("[websocket] unmarshal envelope error: %v", err)
		return
	}
	uids, ok := s.Get("uid")
	if !ok {
		log.Error("[websocket] handleBinaryMessage get uid error, session not contains uid key")
		return
	}
	uid := uids.(int64)

	// 业务消息分发
	g.dispatch(uid, meta)
}

// 错误时调用
func (g *Gate) handleError(s *melody.Session, err error) {
	log.Errorf("[websocket] error occurred, err: %v", err)
}

func (g *Gate) handleClose(s *melody.Session, code int, reason string) error {
	log.Infof("[websocket] connection closed, code: %v, reason: %v", code, reason)
	return nil
}
