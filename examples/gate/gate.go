package main

import (
	"github.com/byteweap/wukong"
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/server/gate"
)

func main() {
	err := wukong.New(
		wukong.ID("gate-1"),
		wukong.Name("gate"),
		wukong.Version("v1.0.0"),
		wukong.Metadata(map[string]string{"author": "Leo"}),
		wukong.Server(gate.New()),
	).Run()
	if err != nil {
		log.Info(err)
	}
}
