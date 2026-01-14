package wats

import (
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Client NATS 和 JetStream 包装器
type Client struct {
	js jetstream.JetStream
}

// Connect 连接到 NATS 服务器并创建 JetStream 客户端
func Connect(urls []string, opts ...Option) (*Client, error) {
	// 处理配置选项
	o := defaultOptions()
	for _, opt := range opts {
		opt.apply(o)
	}

	// 连接 NATS
	url := strings.Join(urls, ",")
	nc, err := nats.Connect(url, o.natsOpts...)
	if err != nil {
		return nil, err
	}

	// 创建 JetStream
	js, err := jetstream.New(nc, o.jsOpts...)
	if err != nil {
		nc.Close()
		return nil, err
	}

	return &Client{js: js}, nil
}

// JetStream 返回 JetStream 实例
func (c *Client) JetStream() jetstream.JetStream {
	return c.js
}

// Nats 返回 NATS 连接
func (c *Client) Nats() *nats.Conn {
	return c.js.Conn()
}

// Close 关闭客户端连接
func (c *Client) Close() {
	conn := c.Nats()
	if conn != nil {
		if err := conn.Drain(); err != nil {
			conn.Close()
		}
	}
}
