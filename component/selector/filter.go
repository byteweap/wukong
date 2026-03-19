package selector

import "context"

// NodeFilter 是节点过滤器
type NodeFilter func(context.Context, []Node) []Node
