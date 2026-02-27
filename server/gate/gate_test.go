package gate_test

import (
	"context"
	"testing"

	"github.com/byteweap/wukong"
	"github.com/byteweap/wukong/server/gate"
	"github.com/stretchr/testify/require"
)

func TestGate(t *testing.T) {

	g := gate.New(
		gate.Websocket(":8080", "/ws"),
		//gate.Locator(locator.Locator()),
		//gate.Broker(broker.Broker("")),
		//gate.Discovery(registry.Discovery("")),
	)
	ctx := wukong.NewContext(context.Background(), wukong.New())

	require.Nil(t, g.Start(ctx))

}
