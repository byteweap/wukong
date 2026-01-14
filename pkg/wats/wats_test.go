package wats_test

import (
	"testing"

	"github.com/nats-io/nats.go"

	"github.com/byteweap/wukong/pkg/wats"
)

func TestKats(t *testing.T) {

	cli, err := wats.Connect(
		[]string{
			"nats://localhost:4222",
			"nats://localhost:4223",
			"nats://localhost:4224",
		},
		wats.WithNats(
			nats.Name("aaaaa"),
		),
		wats.WithJetStream(),
	)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("connect success!")
	cli.Close()
}
