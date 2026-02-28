package gate_test

import (
	"context"
	"testing"
	"time"

	"github.com/byteweap/wukong"
	"github.com/byteweap/wukong/server/gate"
	"github.com/stretchr/testify/require"
)

func TestGate(t *testing.T) {

	g := gate.New(
		gate.SubjectPrefix("leo"),
		gate.Addr(":8080"),
		gate.Path("/"),
		gate.WriteTimeout(time.Second*5),
		gate.PongTimeout(time.Second*60),
		gate.PingInterval(time.Second*5),
		gate.MaxMessageSize(1024*2),
		gate.MessageBufferSize(256),
		//gate.Locator(locator.Locator()),
		//gate.Broker(broker.Broker("")),
		//gate.Discovery(registry.Discovery("")),
	)
	ctx := wukong.NewContext(context.Background(), wukong.New())

	require.Nil(t, g.Start(ctx))

}
