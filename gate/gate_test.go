package gate_test

import (
	"testing"

	"github.com/byteweap/wukong/gate"
	"github.com/stretchr/testify/assert"
)

func TestGate(t *testing.T) {

	g, err := gate.New()
	assert.Nil(t, err)

	t.Cleanup(func() {
		g.Close()
	})

	g.Serve()
}
