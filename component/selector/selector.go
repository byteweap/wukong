package selector

// Selector 选择器接口
type Selector interface {
	// Select 选择节点
	Select(key string, filters ...Filter) (Node, error)
	// Update 更新节点列表
	Update(nodes []Node)
	// Nodes 获取节点列表
	Nodes() []Node
}
