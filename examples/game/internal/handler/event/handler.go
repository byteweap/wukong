package event

import "github.com/byteweap/wukong/examples/game/internal/service"

type EventHandler struct {
	svc service.EventService
}

func New(svc service.EventService) *EventHandler {
	return &EventHandler{
		svc: svc,
	}
}
