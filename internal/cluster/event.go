package cluster

type Event string

const (
	Event_Business  Event = "business"  // 业务 [DEFAULT]
	Event_Online    Event = "online"    // 上线
	Event_Offline   Event = "offline"   // 掉线
	Event_Reconnect Event = "reconnect" // 重连
)
