package main

import (
	"github.com/byteweap/wukong"
	"github.com/byteweap/wukong/component/log"
)

func main() {
	s := wukong.New(
		wukong.ID("gate-1"),
		wukong.Name("gate"),
		wukong.Version("v1.0.0"),
		wukong.Metadata(map[string]string{"author": "Leo"}),
	)
	if err := s.Run(); err != nil {
		log.Errorf("start failed, err: %v", err)
	}
}
