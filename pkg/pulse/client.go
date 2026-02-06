package pulse

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
)

// Client pulse 客户端
// 基于 gobwas/ws 实现
type Client struct {
	opts *clientOptions

	mu     sync.RWMutex
	conn   *ClientConn
	url    string
	header http.Header

	ctx    context.Context
	cancel context.CancelFunc

	closed  atomic.Bool
	running atomic.Bool
}

func NewClient(opts ...ClientOption) *Client {
	o := defaultClientOptions()
	for _, opt := range opts {
		opt(o)
	}
	return &Client{opts: o}
}

func (c *Client) Dial(ctx context.Context, url string, header http.Header) error {
	if c.closed.Load() {
		return ErrClosed
	}
	if !c.running.CompareAndSwap(false, true) {
		return ErrAlreadyRunning
	}

	c.mu.Lock()
	c.url = url
	c.header = header
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.mu.Unlock()

	go c.run()

	return nil
}

func (c *Client) Close() error {
	if !c.closed.CompareAndSwap(false, true) {
		return nil
	}
	c.mu.Lock()
	if c.cancel != nil {
		c.cancel()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	c.mu.Unlock()
	return nil
}

func (c *Client) Conn() *ClientConn {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn
}

func (c *Client) WriteBinary(msg []byte) error {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()
	if conn == nil {
		return ErrNotConnected
	}
	return conn.WriteBinary(msg)
}

func (c *Client) WriteText(msg []byte) error {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()
	if conn == nil {
		return ErrNotConnected
	}
	return conn.WriteText(msg)
}

func (c *Client) run() {
	defer c.running.Store(false)

	backoff := c.opts.reconnectInterval
	for {
		if c.closed.Load() {
			return
		}

		conn, err := c.dialOnce()
		if err != nil {
			c.callError(nil, err)
			if !c.waitReconnect(backoff) {
				return
			}
			backoff = c.nextBackoff(backoff)
			continue
		}

		backoff = c.opts.reconnectInterval
		c.setConn(conn)
		c.callOpen(conn)

		err = conn.run()

		c.setConn(nil)
		if err != nil {
			c.callError(conn, err)
		}

		if c.closed.Load() {
			c.callClose(conn, err)
			return
		}

		if !c.waitReconnect(backoff) {
			c.callClose(conn, err)
			return
		}

		backoff = c.nextBackoff(backoff)
	}
}

func (c *Client) dialOnce() (*ClientConn, error) {
	c.mu.RLock()
	url := c.url
	header := c.header
	ctx := c.ctx
	c.mu.RUnlock()

	dialer := ws.Dialer{Timeout: c.opts.dialTimeout}
	if header != nil {
		dialer.Header = ws.HandshakeHeaderHTTP(header)
	}
	raw, br, _, err := dialer.Dial(ctx, url)
	if err != nil {
		return nil, err
	}

	conn := newClientConn(raw, br, c.opts)
	return conn, nil
}

func (c *Client) waitReconnect(d time.Duration) bool {
	if d <= 0 || c.opts.reconnectInterval <= 0 {
		return false
	}
	if c.ctx == nil {
		return false
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-c.ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

func (c *Client) nextBackoff(current time.Duration) time.Duration {
	if current <= 0 {
		return c.opts.reconnectInterval
	}
	factor := c.opts.reconnectFactor
	if factor <= 1 {
		factor = 2
	}
	next := time.Duration(float64(current) * factor)
	if maxInterval := c.opts.reconnectMaxInterval; maxInterval > 0 && next > maxInterval {
		next = maxInterval
	}
	if next < c.opts.reconnectInterval {
		next = c.opts.reconnectInterval
	}
	return next
}

func (c *Client) setConn(conn *ClientConn) {
	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()
}

func (c *Client) callOpen(conn *ClientConn) {
	if c.opts.onOpen != nil {
		c.opts.onOpen(conn)
	}
}

func (c *Client) callClose(conn *ClientConn, err error) {
	if c.opts.onClose != nil {
		c.opts.onClose(conn, err)
	}
}

func (c *Client) callError(conn *ClientConn, err error) {
	if c.opts.onError != nil {
		c.opts.onError(conn, err)
	}
}
