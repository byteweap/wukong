package gate

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/olahol/melody"

	"github.com/byteweap/meta/component/log"
	"github.com/byteweap/meta/encoding/proto"
	"github.com/byteweap/meta/envelope"
	"github.com/byteweap/meta/internal/cluster"
)

// 连接建立时调用
func (g *Gate) handleConnect(s *melody.Session) {
	req := s.Request

	uid := g.opts.userIDExtractor(req)
	if uid <= 0 {
		_ = s.Write([]byte("uid is required"))
		_ = s.Close()
		return
	}
	// 注册会话
	session, ok := g.sessions.replace(uid, s)
	if ok {
		log.Warnf("[websocket] connection exists: uid: %v, close old connection", uid)
		_ = session.Close()
	}
	s.Set("uid", uid)

	log.Infof("[websocket] new connection success, uid: %v, %s", uid, sessionRemoteAddr(s))

	loc := g.opts.locator

	// 绑定网关
	if err := loc.Bind(g.ctx, uid, g.appName, g.appID); err != nil {
		log.Errorf("[websocket] new connection bind gate error, uid: %v, err: %v", uid, err)
		g.sessions.unregisterIfSame(uid, s)
		_ = s.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		_ = s.CloseWithMsg(websocket.FormatCloseMessage(melody.CloseInternalServerErr, "bind gate error"))
		return
	}

	// 广播 上线、重连 事件到上游服务
	event := cluster.EventOnline
	if ok {
		event = cluster.EventReconnect
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
	if !g.sessions.unregisterIfSame(uid, s) {
		log.Warnf("[websocket] connection disconnect error, uid: %v session not match", uid)
		return
	}

	log.Infof("[websocket] connection disconnect success, uid: %v", uid)

	// 解绑网关
	if err := g.opts.locator.UnBind(g.ctx, uid, g.appName, g.appID); err != nil {
		log.Errorf("[websocket] connection disconnect success, unbind gate error, uid: %v, err: %v", uid, err)
	}

	// 广播掉线事件到上游服务
	g.broadcastEvent(uid, cluster.EventOffline)
}

// 接收到文本消息时调用
func (g *Gate) handleTextMessage(_ *melody.Session, _ []byte) {
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
func (g *Gate) handleError(_ *melody.Session, err error) {
	log.Errorf("[websocket] error occurred, err: %v", err)
}

func (g *Gate) handleClose(_ *melody.Session, code int, reason string) error {
	log.Infof("[websocket] connection closed, code: %v, reason: %v", code, reason)
	return nil
}

func sessionRemoteAddr(s *melody.Session) string {
	if s == nil {
		return ""
	}
	conn := s.WebsocketConnection()
	if conn == nil || conn.RemoteAddr() == nil {
		return ""
	}
	return conn.RemoteAddr().String()
}
