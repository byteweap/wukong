package selector

type Filter func([]Node) []Node

// Version 筛选指定版本节点
func Version(version string) Filter {
	return func(nodes []Node) []Node {
		var res []Node
		for _, n := range nodes {
			if n.Version() == version {
				res = append(res, n)
			}
		}
		return res
	}
}
