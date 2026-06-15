package redis

import (
	"errors"
	"testing"

	goredis "github.com/redis/go-redis/v9"
)

func TestNormalizeNodeResultRedisNil(t *testing.T) {
	node, err := normalizeNodeResult("", goredis.Nil)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if node != "" {
		t.Fatalf("expected empty node, got %q", node)
	}
}

func TestNormalizeNodeResultHit(t *testing.T) {
	node, err := normalizeNodeResult("node-1", nil)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if node != "node-1" {
		t.Fatalf("expected node-1, got %q", node)
	}
}

func TestNormalizeNodeResultError(t *testing.T) {
	wantErr := errors.New("redis unavailable")

	node, err := normalizeNodeResult("", wantErr)

	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
	if node != "" {
		t.Fatalf("expected empty node, got %q", node)
	}
}
