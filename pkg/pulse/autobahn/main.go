package main

import (
	"log"
	"net/http"
	"time"

	"github.com/byteweap/wukong/pkg/pulse"
)

func main() {
	srv := pulse.New(
		pulse.SendQueueSize(256),
		pulse.CheckOrigin(func(origin string) bool { return true }),
		pulse.MaxMessageSize(16*1024*1024),
		pulse.ReadTimeout(time.Second*5),
		pulse.WriteTimeout(time.Second*5),
		pulse.Backpressure(pulse.BackpressureKick),
		pulse.OnConnect(func(c *pulse.Conn) {
			log.Printf("new connection: %s", c.RemoteAddr())
		}),
		pulse.OnDisconnect(func(c *pulse.Conn, err error) {
			log.Printf("connection closed: %s, error: %v", c.RemoteAddr(), err)
		}),
		pulse.OnTextMessage(func(c *pulse.Conn, msg []byte) {
			_ = c.WriteText(msg)
		}),
		pulse.OnBinaryMessage(func(c *pulse.Conn, msg []byte) {
			_ = c.WriteBinary(msg)
		}),
		pulse.OnError(func(c *pulse.Conn, err error) {
			log.Printf("connection error: %s, error: %v", c.RemoteAddr(), err)
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
