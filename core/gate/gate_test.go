package gate_test

import (
	"testing"

	"github.com/byteweap/wukong/core/gate"
)

func TestGate(t *testing.T) {
	g := gate.New()
	g.Start()
}
