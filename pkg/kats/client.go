package kats

import (
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Client nats and jetstream wrapper
type Client struct {
	js jetstream.JetStream
}

func Connect(urls []string, opts ...Option) (*Client, error) {
	// options
	o := defaultOptions()
	for _, opt := range opts {
		opt.apply(o)
	}

	// nats connect
	url := strings.Join(urls, ",")
	nc, err := nats.Connect(url, o.natsOpts...)
	if err != nil {
		return nil, err
	}

	// jetstream
	js, err := jetstream.New(nc, o.jsOpts...)
	if err != nil {
		nc.Close()
		return nil, err
	}

	return &Client{js: js}, nil
}

func (c *Client) JetStream() jetstream.JetStream {
	return c.js
}

func (c *Client) Nats() *nats.Conn {
	return c.js.Conn()
}

func (c *Client) Shutdown() {
	conn := c.Nats()
	if conn != nil {
		conn.Drain()
	}
}
