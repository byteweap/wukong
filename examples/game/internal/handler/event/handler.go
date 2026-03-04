package event

import (
	"github.com/byteweap/wukong/examples/game/internal/handler"
)

type EventHandler struct {
	gs handler.IService
}

func New(gs handler.IService) *EventHandler {
	return &EventHandler{gs: gs}
}
