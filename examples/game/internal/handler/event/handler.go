package event

import "github.com/byteweap/wukong/examples/game/internal/service"

// EventHandler 事件处理
type EventHandler struct {
	svc service.EventService
}

func New(svc service.EventService) *EventHandler {
	return &EventHandler{
		svc: svc,
	}
}
