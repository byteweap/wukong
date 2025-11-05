package websocket

import (
	"context"
	"log"
	"testing"

	"github.com/gobwas/ws"
)

func TestClient(t *testing.T) {

	conn, _, _, err := ws.Dial(context.Background(), "ws://localhost:8089/ws")
	if err != nil {
		log.Printf("[WebSocket] Dial error: %v", err)
		return
	}
	defer conn.Close()

	// 读取消息
	for {
		frame, err := ws.ReadFrame(conn)
		if err != nil {
			log.Printf("[WebSocket] ReadFrame error: %v", err)
			return
		}
		switch frame.Header.OpCode {
		case ws.OpText, ws.OpBinary:
			log.Printf("[WebSocket] Received message: %s", frame.Payload)
		case ws.OpClose:
			log.Printf("[WebSocket] Received close frame")
			//ws.WriteFrame(conn, ws.NewCloseFrame(ws.NewCloseFrameBody(ws.StatusNormalClosure, "")))
			return
		case ws.OpPing:
			log.Printf("[WebSocket] Received ping frame")
			// ws.WriteFrame(conn, ws.NewPongFrame(frame.Payload))
		case ws.OpPong:
			log.Printf("[WebSocket] Received pong frame")
		default:
			log.Printf("[WebSocket] Received unknown frame: %v", frame)
		}
	}
}
