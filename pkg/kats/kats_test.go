package kats_test

import (
	"testing"

	"github.com/nats-io/nats.go"

	"github.com/byteweap/wukong/pkg/kats"
)

func TestKats(t *testing.T) {

	cli, err := kats.Connect(
		[]string{"nats://localhost:4222", "nats://localhost:4223", "nats://localhost:4224"},
		kats.WithNats(
			nats.Name("aaaaa"),
		),
		kats.WithJetStream(),
	)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("connect success!")
	cli.Shutdown()
}
