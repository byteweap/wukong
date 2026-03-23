package random

import (
	"errors"
	"testing"

	"github.com/byteweap/wukong/component/selector"
)

func TestRandomSelector_Select_EmptyNodes(t *testing.T) {
	rs := NewRandomSelector()
	rs.Update([]selector.Node{})

	got, err := rs.Select("k")
	if !errors.Is(err, selector.ErrNoAvailableNode) {
		t.Fatalf("expected ErrNoAvailableNode, got %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil node for empty nodes, got %#v", got)
	}
}

func TestRandomSelector_UpdateAndNodes(t *testing.T) {
	rs := NewRandomSelector()
	nodes := []selector.Node{
		selector.NewNode("node-a", "service1", "mesh", "v1.0.0", 1, map[string]any{"zone": "1"}),
		selector.NewNode("node-b", "service1", "mesh", "v1.0.0", 2, map[string]any{"zone": "2"}),
	}

	rs.Update(nodes)
	got := rs.Nodes()
	if len(got) != len(nodes) || got[0].ID() != "node-a" || got[1].ID() != "node-b" {
		t.Fatalf("Nodes() mismatch: got %#v want %#v", got, nodes)
	}
}

func TestRandomSelector_Select_WeightedBias(t *testing.T) {
	rs := NewRandomSelector()
	rs.Update([]selector.Node{
		selector.NewNode("node-a", "service1", "mesh", "v1.0.0", 2, nil),
		selector.NewNode("node-b", "service1", "mesh", "v1.0.0", 1, nil),
	})

	counts := map[string]int{}
	for i := 0; i < 1000; i++ {
		n, err := rs.Select("k")
		if err != nil {
			t.Fatalf("Select returned error: %v", err)
		}
		counts[n.ID()]++
	}

	if counts["node-a"] <= counts["node-b"] {
		t.Fatalf("expected weighted bias towards node-a, got %v", counts)
	}
}

func TestRandomSelector_Update_DefaultWeightForNonPositive(t *testing.T) {
	rs := NewRandomSelector()
	rs.Update([]selector.Node{
		selector.NewNode("node-a", "service1", "mesh", "v1.0.0", 0, nil),
		selector.NewNode("node-b", "service1", "mesh", "v1.0.0", 1, nil),
	})

	if rs.totalWeight != 2 {
		t.Fatalf("expected totalWeight 2, got %v", rs.totalWeight)
	}
	if len(rs.cumulative) != 2 || rs.cumulative[0] != 1 || rs.cumulative[1] != 2 {
		t.Fatalf("unexpected cumulative weights: %v", rs.cumulative)
	}
}

func TestRandomSelector_Select_WithFilter(t *testing.T) {
	rs := NewRandomSelector()
	rs.Update([]selector.Node{
		selector.NewNode("node-a", "service1", "mesh", "v1.0.0", 1, nil),
		selector.NewNode("node-b", "service1", "mesh", "v1.0.0", 1, nil),
	})

	onlyB := func(nodes []selector.Node) []selector.Node {
		out := make([]selector.Node, 0, len(nodes))
		for _, n := range nodes {
			if n.ID() == "node-b" {
				out = append(out, n)
			}
		}
		return out
	}

	got, err := rs.Select("k", onlyB)
	if err != nil {
		t.Fatalf("Select returned error: %v", err)
	}
	if got.ID() != "node-b" {
		t.Fatalf("expected node-b, got %v", got.ID())
	}
}
