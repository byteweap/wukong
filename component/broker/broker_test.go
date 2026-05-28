package broker

import (
	"testing"
)

func TestHeader_Add(t *testing.T) {
	h := Header{}
	h.Add("key", "value1")
	h.Add("key", "value2")

	values := h.Values("key")
	if len(values) != 2 {
		t.Errorf("expected 2 values, got %d", len(values))
	}
	if values[0] != "value1" {
		t.Errorf("expected 'value1', got '%s'", values[0])
	}
	if values[1] != "value2" {
		t.Errorf("expected 'value2', got '%s'", values[1])
	}
}

func TestHeader_Set(t *testing.T) {
	h := Header{}
	h.Add("key", "value1")
	h.Set("key", "value2")

	values := h.Values("key")
	if len(values) != 1 {
		t.Errorf("expected 1 value, got %d", len(values))
	}
	if values[0] != "value2" {
		t.Errorf("expected 'value2', got '%s'", values[0])
	}
}

func TestHeader_Get(t *testing.T) {
	t.Run("nil header", func(t *testing.T) {
		var h Header
		if h.Get("key") != "" {
			t.Error("expected empty string for nil header")
		}
	})

	t.Run("existing key", func(t *testing.T) {
		h := Header{}
		h.Set("key", "value")
		if h.Get("key") != "value" {
			t.Errorf("expected 'value', got '%s'", h.Get("key"))
		}
	})

	t.Run("missing key", func(t *testing.T) {
		h := Header{}
		if h.Get("key") != "" {
			t.Error("expected empty string for missing key")
		}
	})
}

func TestHeader_Values(t *testing.T) {
	h := Header{}
	h.Add("key", "value1")
	h.Add("key", "value2")

	values := h.Values("key")
	if len(values) != 2 {
		t.Errorf("expected 2 values, got %d", len(values))
	}

	missing := h.Values("missing")
	if missing != nil {
		t.Errorf("expected nil for missing key, got %v", missing)
	}
}

func TestHeader_Del(t *testing.T) {
	h := Header{}
	h.Set("key", "value")
	h.Del("key")

	if h.Get("key") != "" {
		t.Error("expected empty string after deletion")
	}
}

func TestPubHeader(t *testing.T) {
	h := Header{"key": {"value"}}
	opt := PubHeader(h)

	opts := &PublishOptions{}
	opt(opts)

	if opts.Header == nil {
		t.Error("expected header to be set")
	}
	if opts.Header.Get("key") != "value" {
		t.Errorf("expected 'value', got '%s'", opts.Header.Get("key"))
	}
}

func TestRequestHeader(t *testing.T) {
	h := Header{"key": {"value"}}
	opt := RequestHeader(h)

	opts := &RequestOptions{}
	opt(opts)

	if opts.Header == nil {
		t.Error("expected header to be set")
	}
	if opts.Header.Get("key") != "value" {
		t.Errorf("expected 'value', got '%s'", opts.Header.Get("key"))
	}
}

func TestSubQueue(t *testing.T) {
	opt := SubQueue("queue1")

	opts := &SubscribeOptions{}
	opt(opts)

	if opts.Queue != "queue1" {
		t.Errorf("expected 'queue1', got '%s'", opts.Queue)
	}
}

func TestReplyHeader(t *testing.T) {
	h := Header{"key": {"value"}}
	opt := ReplyHeader(h)

	opts := &ReplyOptions{}
	opt(opts)

	if opts.Header == nil {
		t.Error("expected header to be set")
	}
	if opts.Header.Get("key") != "value" {
		t.Errorf("expected 'value', got '%s'", opts.Header.Get("key"))
	}
}

func TestPubReply(t *testing.T) {
	opt := PubReply("reply.subject")

	opts := &PublishOptions{}
	opt(opts)

	if opts.Reply != "reply.subject" {
		t.Errorf("expected 'reply.subject', got '%s'", opts.Reply)
	}
}

func TestMessage(t *testing.T) {
	msg := &Message{
		Subject: "test.subject",
		Reply:   "test.reply",
		Header:  Header{"key": {"value"}},
		Data:    []byte("test data"),
	}

	if msg.Subject != "test.subject" {
		t.Errorf("expected 'test.subject', got '%s'", msg.Subject)
	}
	if msg.Reply != "test.reply" {
		t.Errorf("expected 'test.reply', got '%s'", msg.Reply)
	}
	if msg.Header.Get("key") != "value" {
		t.Errorf("expected 'value', got '%s'", msg.Header.Get("key"))
	}
	if string(msg.Data) != "test data" {
		t.Errorf("expected 'test data', got '%s'", string(msg.Data))
	}
}
