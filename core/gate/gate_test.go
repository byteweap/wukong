package gate_test

import (
	"testing"

	"github.com/byteweap/wukong/core/gate"
	"github.com/byteweap/wukong/pkg/klog"
)

func TestGate(t *testing.T) {
	logger := klog.New()
	g := gate.New(logger)
	g.Start()
}
