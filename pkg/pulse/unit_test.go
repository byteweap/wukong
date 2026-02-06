package pulse

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/gobwas/ws"
)

func TestServerBackpressureKick(t *testing.T) {
	c1, c2 := net.Pipe()
	defer func() {
		_ = c1.Close()
		_ = c2.Close()
	}()

	opts := &options{
		sendQueueSize: 1,
		backpressure:  BackpressureKick,
	}

	conn := &Conn{
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
	defer func() {
		_ = c1.Close()
		_ = c2.Close()
	}()

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

func TestServerCloseSendsCloseFrameWhenQueueFull(t *testing.T) {
	c1, c2 := net.Pipe()
	defer func() {
		_ = c1.Close()
		_ = c2.Close()
	}()

	opts := &options{
		sendQueueSize: 1,
		backpressure:  BackpressureBlock,
	}

	conn := &Conn{
		raw:   c1,
		opts:  opts,
		sendQ: make(chan serverSendItem, opts.sendQueueSize),
		done:  make(chan struct{}),
	}

	conn.sendQ <- serverSendItem{op: ws.OpText, msg: []byte("full")}

	conn.Close()

	_ = c2.SetReadDeadline(time.Now().Add(1 * time.Second))
	frame, err := ws.ReadFrame(c2)
	if err != nil {
		t.Fatalf("read frame failed: %v", err)
	}
	if frame.Header.OpCode != ws.OpClose {
		t.Fatalf("expected close frame, got %v", frame.Header.OpCode)
	}
}

func TestClientCloseSendsCloseFrameWhenQueueFull(t *testing.T) {
	c1, c2 := net.Pipe()
	defer func() {
		_ = c2.Close()
	}()

	opts := &clientOptions{
		sendQueueSize: 1,
		backpressure:  BackpressureBlock,
	}

	conn := &ClientConn{
		raw:   c1,
		opts:  opts,
		sendQ: make(chan clientSendItem, opts.sendQueueSize),
		done:  make(chan struct{}),
	}

	conn.sendQ <- clientSendItem{op: ws.OpText, msg: []byte("full")}

	conn.Close()

	_ = c2.SetReadDeadline(time.Now().Add(1 * time.Second))
	frame, err := ws.ReadFrame(c2)
	if err != nil {
		t.Fatalf("read frame failed: %v", err)
	}
	if frame.Header.OpCode != ws.OpClose {
		t.Fatalf("expected close frame, got %v", frame.Header.OpCode)
	}
	if !frame.Header.Masked {
		t.Fatalf("expected masked close frame from client")
	}
}
