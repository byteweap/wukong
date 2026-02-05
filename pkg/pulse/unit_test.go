package pulse

import (
	"errors"
	"net"
	"testing"
	"time"
)

func TestServerBackpressureKick(t *testing.T) {
	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	opts := &serverOptions{
		sendQueueSize: 1,
		backpressure:  BackpressureKick,
	}

	conn := &ServerConn{
		raw:   c1,
		opts:  opts,
		sendQ: make(chan serverSendItem, opts.sendQueueSize),
		done:  make(chan struct{}),
	}

	if err := conn.WriteBinary([]byte("a")); err != nil {
		t.Fatalf("first write failed: %v", err)
	}
	if err := conn.WriteBinary([]byte("b")); !errors.Is(err, ErrBackpressure) {
		t.Fatalf("expected ErrBackpressure, got %v", err)
	}
	if !conn.closed.Load() {
		t.Fatalf("expected conn closed after backpressure kick")
	}
}

func TestClientBackpressureKick(t *testing.T) {
	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	opts := &clientOptions{
		sendQueueSize: 1,
		backpressure:  BackpressureKick,
	}

	conn := &ClientConn{
		raw:   c1,
		opts:  opts,
		sendQ: make(chan clientSendItem, opts.sendQueueSize),
		done:  make(chan struct{}),
	}

	if err := conn.WriteBinary([]byte("a")); err != nil {
		t.Fatalf("first write failed: %v", err)
	}
	if err := conn.WriteBinary([]byte("b")); !errors.Is(err, ErrBackpressure) {
		t.Fatalf("expected ErrBackpressure, got %v", err)
	}
	if !conn.closed.Load() {
		t.Fatalf("expected conn closed after backpressure kick")
	}
}

func TestClientNextBackoff(t *testing.T) {
	c := NewClient(
		ClientReconnectInterval(100*time.Millisecond),
		ClientReconnectMaxInterval(250*time.Millisecond),
		ClientReconnectFactor(2),
	)

	if got := c.nextBackoff(0); got != 100*time.Millisecond {
		t.Fatalf("expected 100ms, got %v", got)
	}
	if got := c.nextBackoff(100 * time.Millisecond); got != 200*time.Millisecond {
		t.Fatalf("expected 200ms, got %v", got)
	}
	if got := c.nextBackoff(200 * time.Millisecond); got != 250*time.Millisecond {
		t.Fatalf("expected 250ms (capped), got %v", got)
	}
}
