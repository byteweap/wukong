package event

import (
	"github.com/byteweap/wukong/examples/game/internal/handler"
)

type EventHandler struct {
	gs handler.GameService
}

func New(gs handler.GameService) *EventHandler {
	return &EventHandler{gs: gs}
}
