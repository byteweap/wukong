package websocket_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/byteweap/wukong/pkg/knet/websocket"
)

func TestServer_HandleRequest(t *testing.T) {

	ws := websocket.NewServer(
		websocket.WithAddr(":8020"),
	)
	http.HandleFunc("/ws", ws.HandleRequest)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		t.Fatal(err)
	}
}

func TestServer_Run(t *testing.T) {

	ws := websocket.NewServer(
		websocket.WithAddr(":8020"),
		websocket.WithPattern("/ws"),
	)
	go func() {
		ws.Run()
	}()

	time.Sleep(time.Second * 10)

	if err := ws.Shutdown(); err != nil {
		t.Fatal(err)
	}
}
