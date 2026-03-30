package selector

import "github.com/byteweap/meta/pkg/conv"

type Node interface {
	ID() string
	Service() string
	Weight() float64
	Version() string
	Meta() map[string]string
}

type node struct {
	id      string
	service string
	version string
	meta    map[string]string
}

func NewNode(id, service, version string, meta map[string]string) Node {
	return &node{
		id:      id,
		service: service,
		version: version,
		meta:    meta,
	}
}

func (n *node) ID() string {
	return n.id
}
func (n *node) Service() string {
	return n.service
}
func (n *node) Weight() float64 {
	return conv.Float64(n.meta["weight"])
}

func (n *node) Version() string {
	return n.version
}

func (n *node) Meta() map[string]string {
	return n.meta
}
