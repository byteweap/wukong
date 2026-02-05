package main

import (
	"log"
	"net/http"

	"github.com/byteweap/wukong/pkg/pulse"
	"github.com/gobwas/ws"
)

func main() {
	srv := pulse.NewServer(
		pulse.ServerMaxMessageSize(16*1024*1024),
		pulse.OnServerMessage(func(c *pulse.ServerConn, op ws.OpCode, msg []byte) {
			switch op {
			case ws.OpText:
				_ = c.WriteText(msg)
			case ws.OpBinary:
				_ = c.WriteBinary(msg)
			}
		}),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if err := srv.HandleRequest(w, r); err != nil {
			log.Printf("upgrade error: %v", err)
			return
		}
	})

	addr := "127.0.0.1:9001"
	log.Printf("autobahn ws server listening on %s/ws", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("listen failed: %v", err)
	}
}
