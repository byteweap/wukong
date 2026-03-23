package selector

type Node interface {
	ID() string
	App() string
	Weight() float64
	Scheme() string
	Version() string
	Meta() map[string]any
}

type node struct {
	id      string
	app     string
	weight  float64
	scheme  string
	version string
	meta    map[string]any
}

func NewNode(id, app, scheme, version string, weight float64, meta map[string]any) Node {
	return &node{
		id:      id,
		app:     app,
		weight:  weight,
		scheme:  scheme,
		version: version,
		meta:    meta,
	}
}

func (n *node) ID() string {
	return n.id
}
func (n *node) App() string {
	return n.app
}
func (n *node) Weight() float64 {
	return n.weight
}
func (n *node) Scheme() string {
	return n.scheme
}

func (n *node) Version() string {
	return n.version
}
func (n *node) Meta() map[string]any {
	return n.meta
}
