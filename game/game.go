package game

import (
	"context"

	"github.com/byteweap/wukong/component/log"
)

// Game 游戏逻辑服
type Game struct {
	opts   *options
	ctx    context.Context
	cancel context.CancelFunc
}

func New(opts ...Option) *Game {
	o := defaultOptions()

	for _, opt := range opts {
		opt(o)
	}
	if o.logger != nil {
		log.SetLogger(o.logger)
	}
	return &Game{opts: o}
}

// Run 运行游戏
func (g *Game) Run() error {
	// todo
	return nil
}

// Stop 停止游戏
func (g *Game) Stop() error {
	// todo
	return nil
}
