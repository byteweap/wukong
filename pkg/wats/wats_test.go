package wats

import (
	"testing"

	"github.com/nats-io/nats.go"
)

func TestWats(t *testing.T) {

	cli, err := Connect(
		[]string{"nats://localhost:4222", "nats://localhost:4223", "nats://localhost:4224"},
		WithNats(
			nats.Name("aaaaa"),
		),
		WithJetStream(),
	)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("connect success!")
	cli.Shutdown()
}
