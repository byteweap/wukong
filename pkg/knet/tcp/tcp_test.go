package tcp

import (
	"math/rand"
	"net"
	"testing"
	"time"
)

func TestTcpServer(t *testing.T) {

	s := NewServer()
	s.Start()
}

func TestTcpClient(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	go func() {

		for {
			num := rand.Intn(10)
			time.Sleep(time.Millisecond * time.Duration(num) * 100)
			conn.Write([]byte("hello world"))
		}
	}()

	go func() {
		for {
			buf := make([]byte, 1024)
			cnt, err := conn.Read(buf)
			if err != nil {
				t.Error(err)
				return
			}
			t.Logf("receive msg: %s", buf[:cnt])
		}
	}()

	select {}
}
