package main

import (
	"github.com/byteweap/wukong/core/gate"
)

//
//func main() {
//	s := websocket.NewServer()
//
//	http.HandleFunc("/ws", s.HandleRequest)
//
//	go func() {
//
//		time.Sleep(time.Second * 5)
//		conn, _, _, err := ws.Dial(context.Background(), "ws://localhost:8089/ws")
//		if err != nil {
//			log.Printf("[WebSocket] Dial error: %v", err)
//			return
//		}
//		defer conn.Close()
//
//		// 读取消息
//		for {
//			frame, err := ws.ReadFrame(conn)
//			if err != nil {
//				log.Printf("[WebSocket] ReadFrame error: %v", err)
//				return
//			}
//			switch frame.Header.OpCode {
//			case ws.OpText, ws.OpBinary:
//				log.Printf("[WebSocket] Received message: %s", frame.Payload)
//			case ws.OpClose:
//				log.Printf("[WebSocket] Received close frame")
//				ws.WriteFrame(conn, ws.NewCloseFrame(ws.NewCloseFrameBody(ws.StatusNormalClosure, "")))
//				return
//			case ws.OpPing:
//				log.Printf("[WebSocket] Received ping frame")
//				// ws.WriteFrame(conn, ws.NewPongFrame(frame.Payload))
//			case ws.OpPong:
//				log.Printf("[WebSocket] Received pong frame")
//			default:
//				log.Printf("[WebSocket] Received unknown frame: %v", frame)
//			}
//		}
//	}()
//
//	http.ListenAndServe(":8089", nil)
//
//}

func main() {

	g := gate.New()
	g.Start()
}
