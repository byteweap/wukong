package gate_test

import (
	"testing"

	"github.com/byteweap/wukong/contrib/logger/zerolog"
	"github.com/byteweap/wukong/gate"
)

func TestGate(t *testing.T) {

	logger := zerolog.New()

	g := gate.New(logger)
	g.Start()
}
