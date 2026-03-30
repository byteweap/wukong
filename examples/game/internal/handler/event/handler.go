package event

import "github.com/byteweap/meta/examples/game/internal/service"

// EventHandler 事件处理
type EventHandler struct {
	svc service.EventService
}

func New(svc service.EventService) *EventHandler {
	return &EventHandler{
		svc: svc,
	}
}
