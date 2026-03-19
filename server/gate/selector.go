package gate

import (
	"github.com/byteweap/wukong/component/registry"
	"github.com/byteweap/wukong/component/selector"
	defaultselector "github.com/byteweap/wukong/contrib/selector"
)

func (g *Gate) selectorFor(service string) selector.Selector {
	if service == "" {
		return g.opts.selector.Build()
	}
	if v, ok := g.selectors.Load(service); ok {
		return v.(selector.Selector)
	}
	sel := g.opts.selector.Build()
	actual, _ := g.selectors.LoadOrStore(service, sel)
	return actual.(selector.Selector)
}

func (g *Gate) buildSelectorNodes(services []*registry.ServiceInstance) []selector.Node {
	nodes := make([]selector.Node, 0, len(services))
	for _, ins := range services {
		if ins == nil || ins.ID == "" {
			continue
		}
		nodes = append(nodes, defaultselector.NewNode("", ins.ID, ins))
	}
	return nodes
}
