package event

import "github.com/byteweap/wukong/examples/game/internal/game"

type EventHandler struct {
	game *game.Game
}

func New(g *game.Game) *EventHandler {
	return &EventHandler{game: g}
}
