package pulse

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gobwas/ws"
)

func TestClientServerMessageFlow(t *testing.T) {
	recvServer := make(chan string, 1)
	recvClient := make(chan string, 1)

	srv := New(
		OnTextMessage(func(c *Conn, msg []byte) {
			recvServer <- string(msg)
			_ = c.WriteText([]byte("world"))
		}),
	)

	h := http.NewServeMux()
	h.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if err := srv.HandleRequest(w, r); err != nil {
			return
		}
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

	client := NewClient(
		OnClientOpen(func(c *ClientConn) {
			_ = c.WriteText([]byte("hello"))
		}),
		OnClientMessage(func(_ *ClientConn, _ ws.OpCode, msg []byte) {
			recvClient <- string(msg)
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := client.Dial(ctx, wsURL, nil); err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer client.Close()

	select {
	case got := <-recvServer:
		if got != "hello" {
			t.Fatalf("server got %q", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting server message")
	}

	select {
	case got := <-recvClient:
		if got != "world" {
			t.Fatalf("client got %q", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting client message")
	}
}

func TestClientAutoReconnect(t *testing.T) {
	var serverConnCount atomic.Int64
	var clientOpenCount atomic.Int64

	srv := New(
		OnConnect(func(c *Conn) {
			if serverConnCount.Add(1) == 1 {
				go func() {
					time.Sleep(50 * time.Millisecond)
					c.Close()
				}()
			}
		}),
	)

	h := http.NewServeMux()
	h.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		_ = srv.HandleRequest(w, r)
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

	client := NewClient(
		ClientReconnectInterval(50*time.Millisecond),
		ClientReconnectMaxInterval(100*time.Millisecond),
		OnClientOpen(func(_ *ClientConn) {
			clientOpenCount.Add(1)
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := client.Dial(ctx, wsURL, nil); err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer client.Close()

	deadline := time.After(3 * time.Second)
	for {
		if clientOpenCount.Load() >= 2 {
			break
		}
		select {
		case <-deadline:
			t.Fatalf("expected reconnect with 2 opens, got %d", clientOpenCount.Load())
		case <-time.After(20 * time.Millisecond):
		}
	}
}
