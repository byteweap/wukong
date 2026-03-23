package selector

type Node interface {
	ID() string
	App() string
	Weight() float64
	Scheme() string
	Version() string
	Meta() map[string]any
}
