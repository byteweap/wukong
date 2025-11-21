package redis

import "github.com/byteweap/wukong/plugin/locator"

// Locator TODO
type Locator struct {
}

var _ locator.Locator = (*Locator)(nil)

func (l *Locator) ID() string {
	//TODO implement me
	panic("implement me")
}

func (l *Locator) Gate(uid int64) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Locator) BindGate(uid int64, node string) error {
	//TODO implement me
	panic("implement me")
}

func (l *Locator) UnBindGate(uid int64, node string) error {
	//TODO implement me
	panic("implement me")
}

func (l *Locator) Game(uid int64) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Locator) BindGame(uid int64, node string) error {
	//TODO implement me
	panic("implement me")
}

func (l *Locator) UnBindGame(uid int64, node string) error {
	//TODO implement me
	panic("implement me")
}
