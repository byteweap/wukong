package main

import "github.com/byteweap/wukong/gate"

func main() {

	g, err := gate.New()
	if err != nil {
		panic(err)
	}
	defer g.Close()

	g.Serve()

}
