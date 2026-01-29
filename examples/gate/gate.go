package main

import (
	"github.com/byteweap/wukong"
	"github.com/byteweap/wukong/component/log"
)

func main() {
	s := wukong.New()
	if err := s.Run(); err != nil {
		log.Errorf("start failed, err: %v", err)
	}
}
