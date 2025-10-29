package websocket_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/byteweap/wukong/pkg/wnet/websocket"
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
	ws.OnStart(func(addr, pattern string) {
		t.Logf("server start success, addr: %s, pattern: %s", addr, pattern)
	})
	ws.OnStop(func() {
		t.Logf("server stop success")
	})
	ws.OnConnect(func(conn *websocket.Conn) {
		t.Logf("connect success")
	})
	ws.OnMessage(func(conn *websocket.Conn, msg []byte) {
		t.Logf("recieve message success, len: %d", len(msg))
	})
	ws.ErrorHandler(func(err error) {
		t.Logf("websocket error: %s", err.Error())
	})

	go func() {
		ws.Run()
	}()

	time.Sleep(time.Second * 10)
	ws.Shutdown()

}
