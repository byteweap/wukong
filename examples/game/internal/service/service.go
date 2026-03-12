package service

// IRoomService 提供 handler 需要的房间领域能力
// 接口保持精简，按用例逐步扩展
type IRoomService interface {
	NumRooms() int
}

type IPlayerService interface {
	NumPlayers() int
}

// EventService 是事件处理器所需的能力集合
type EventService interface {
	IRoomService
	IPlayerService
}
