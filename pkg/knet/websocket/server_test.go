package websocket_test

import (
	"net/http"
	"testing"

	"github.com/byteweap/wukong/pkg/knet"
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
	ws.OnStart(func(addr, pattern string) {
		t.Logf("server start success, addr: %s, pattern: %s", addr, pattern)
	})
	ws.OnStop(func() {
		t.Logf("server stop success")
	})

	ws.OnConnect(func(conn knet.Conn) {
		t.Logf("connect success, id: %v", conn.ID())
	})
	ws.OnTextMessage(func(conn knet.Conn, msg []byte) {
		t.Logf("recieve text message success, len: %d", len(msg))
	})
	ws.OnBinaryMessage(func(conn knet.Conn, msg []byte) {
		t.Logf("recieve binary message success, len: %d", len(msg))
	})
	ws.OnError(func(err error) {
		t.Logf("websocket error: %s", err.Error())
	})

	ws.Start()

}
