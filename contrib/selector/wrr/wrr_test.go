package wrr

import (
	"errors"
	"reflect"
	"testing"

	"github.com/byteweap/wukong/component/selector"
)

func TestWRRSelector_Select_EmptyNodes(t *testing.T) {
	ws := NewWRRSelector()
	ws.Update([]selector.Node{})

	got, err := ws.Select("ignored")
	if !errors.Is(err, selector.ErrNoAvailableNode) {
		t.Fatalf("expected ErrNoAvailableNode, got %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil node for empty nodes, got %#v", got)
	}
}

func TestWRRSelector_UpdateAndNodes(t *testing.T) {
	ws := NewWRRSelector()
	nodes := []selector.Node{
		selector.NewNode("node-a", "service1", "mesh", "v1.0.0", 1, map[string]any{"zone": "1"}),
		selector.NewNode("node-b", "service1", "mesh", "v1.0.0", 2, map[string]any{"zone": "2"}),
	}

	ws.Update(nodes)
	got := ws.Nodes()
	if len(got) != len(nodes) || got[0].ID() != "node-a" || got[1].ID() != "node-b" {
		t.Fatalf("Nodes() mismatch: got %#v want %#v", got, nodes)
	}
}

func TestWRRSelector_Select_SmoothSequence(t *testing.T) {
	ws := NewWRRSelector()
	ws.Update([]selector.Node{
		selector.NewNode("node-a", "service1", "mesh", "v1.0.0", 2, nil),
		selector.NewNode("node-b", "service1", "mesh", "v1.0.0", 1, nil),
	})

	want := []string{"node-a", "node-b", "node-a", "node-a", "node-b", "node-a"}
	got := make([]string, 0, len(want))
	for range want {
		n, err := ws.Select("ignored")
		if err != nil {
			t.Fatalf("Select returned error: %v", err)
		}
		got = append(got, n.ID())
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sequence mismatch: got %v want %v", got, want)
	}
}

func TestWRRSelector_UpdateResetsState(t *testing.T) {
	ws := NewWRRSelector()
	ws.Update([]selector.Node{
		selector.NewNode("node-a", "service1", "mesh", "v1.0.0", 2, nil),
		selector.NewNode("node-b", "service1", "mesh", "v1.0.0", 1, nil),
	})

	for i := 0; i < 3; i++ {
		if _, err := ws.Select("ignored"); err != nil {
			t.Fatalf("Select returned error: %v", err)
		}
	}

	ws.Update([]selector.Node{
		selector.NewNode("node-c", "service1", "mesh", "v1.0.0", 1, nil),
		selector.NewNode("node-d", "service1", "mesh", "v1.0.0", 1, nil),
	})

	want := []string{"node-c", "node-d", "node-c", "node-d"}
	got := make([]string, 0, len(want))
	for range want {
		n, err := ws.Select("ignored")
		if err != nil {
			t.Fatalf("Select returned error: %v", err)
		}
		got = append(got, n.ID())
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sequence mismatch after Update: got %v want %v", got, want)
	}
}

func TestWRRSelector_Select_DefaultWeightForNonPositive(t *testing.T) {
	ws := NewWRRSelector()
	ws.Update([]selector.Node{
		selector.NewNode("node-a", "service1", "mesh", "v1.0.0", -2, nil),
		selector.NewNode("node-b", "service1", "mesh", "v1.0.0", 1, nil),
	})

	want := []string{"node-a", "node-b", "node-a", "node-b"}
	got := make([]string, 0, len(want))
	for range want {
		n, err := ws.Select("ignored")
		if err != nil {
			t.Fatalf("Select returned error: %v", err)
		}
		got = append(got, n.ID())
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sequence mismatch with non-positive weight: got %v want %v", got, want)
	}
}

func TestWRRSelector_Select_WithFilter(t *testing.T) {
	ws := NewWRRSelector()
	ws.Update([]selector.Node{
		selector.NewNode("node-a", "service1", "mesh", "v1.0.0", 1, nil),
		selector.NewNode("node-b", "service1", "mesh", "v1.0.0", 1, nil),
		selector.NewNode("node-c", "service1", "mesh", "v1.0.0", 1, nil),
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

	got, err := ws.Select("ignored", onlyB)
	if err != nil {
		t.Fatalf("Select returned error: %v", err)
	}
	if got.ID() != "node-b" {
		t.Fatalf("expected node-b, got %v", got.ID())
	}
}
