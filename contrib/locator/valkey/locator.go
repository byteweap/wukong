package valkey

import (
	"context"

	"github.com/byteweap/wukong/component/locator"
)

type Locator struct {
	// Add fields here
}

var _ locator.Locator = (*Locator)(nil)

func (l *Locator) ID() string {
	//TODO implement me
	panic("implement me")
}

func (l *Locator) Gate(ctx context.Context, uid int64) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Locator) BindGate(ctx context.Context, uid int64, node string) error {
	//TODO implement me
	panic("implement me")
}

func (l *Locator) UnBindGate(ctx context.Context, uid int64, node string) error {
	//TODO implement me
	panic("implement me")
}

func (l *Locator) Game(ctx context.Context, uid int64) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Locator) BindGame(ctx context.Context, uid int64, node string) error {
	//TODO implement me
	panic("implement me")
}

func (l *Locator) UnBindGame(ctx context.Context, uid int64, node string) error {
	//TODO implement me
	panic("implement me")
}
