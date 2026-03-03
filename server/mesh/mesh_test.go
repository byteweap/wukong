package mesh_test

import (
	"context"
	"testing"

	"github.com/byteweap/wukong/server/mesh"
)

type Params struct {
	Name string
}

func EnterGame(ctx *mesh.Context, req *Params) error {
	return nil
}

func TestMesh(t *testing.T) {

	app := mesh.New()
	app.Route(1, 1, mesh.Wrap(EnterGame))
	app.Start(context.Background())
}
