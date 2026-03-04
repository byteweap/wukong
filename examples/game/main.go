package main

import (
	"github.com/byteweap/wukong"
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/examples/game/internal"
)

func main() {

	s := internal.New()

	err := wukong.New(
		wukong.ID("game-1"),
		wukong.Name("game"),
		wukong.Version("v1.0.0"),
		wukong.Metadata(map[string]string{"author": "Leo"}),
		wukong.Server(s),
	).Run()
	if err != nil {
		log.Info(err)
	}
}
