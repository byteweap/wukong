// Package filter 提供节点过滤器
package filter

import (
	"context"

	"github.com/byteweap/wukong/component/selector"
)

// Version 返回按版本过滤器
func Version(version string) selector.NodeFilter {
	return func(_ context.Context, nodes []selector.Node) []selector.Node {
		newNodes := make([]selector.Node, 0, len(nodes))
		for _, n := range nodes {
			if n.Version() == version {
				newNodes = append(newNodes, n)
			}
		}
		return newNodes
	}
}
