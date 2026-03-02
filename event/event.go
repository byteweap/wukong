package event

type Event int8

const (
	EventOnline    Event = iota // 上线
	EventOffline                // 掉线
	EventReconnect              // 重连
	EventBussiness              // 业务消息
)
