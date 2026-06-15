package cluster

type Event string

const (
	EventBusiness  Event = "business"  // 业务 [DEFAULT]
	EventOnline    Event = "online"    // 上线
	EventOffline   Event = "offline"   // 掉线
	EventReconnect Event = "reconnect" // 重连
)
