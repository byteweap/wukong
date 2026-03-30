package roundrobin

import (
	"errors"
	"reflect"
	"testing"

	"github.com/byteweap/meta/component/selector"
)

func TestRoundRobinSelector_Select_EmptyNodes(t *testing.T) {
	rr := New()
	rr.Update([]selector.Node{})

	got, err := rr.Select("k")
	if !errors.Is(err, selector.ErrNoAvailableNode) {
		t.Fatalf("expected ErrNoAvailableNode, got %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil node for empty nodes, got %#v", got)
	}
}

func TestRoundRobinSelector_UpdateAndNodes(t *testing.T) {
	rr := New()
	nodes := []selector.Node{
		selector.NewNode("node-a", "service1", "v1.0.0", map[string]string{"weight": "1"}),
		selector.NewNode("node-b", "service1", "v1.0.0", map[string]string{"weight": "2"}),
	}

	rr.Update(nodes)
	got := rr.Nodes()
	if len(got) != len(nodes) || got[0].ID() != "node-a" || got[1].ID() != "node-b" {
		t.Fatalf("Nodes() mismatch: got %#v want %#v", got, nodes)
	}
}

func TestRoundRobinSelector_Select_Sequence(t *testing.T) {
	rr := New()
	rr.Update([]selector.Node{
		selector.NewNode("node-a", "service1", "v1.0.0", map[string]string{"weight": "0"}),
		selector.NewNode("node-b", "service1", "v1.0.0", map[string]string{"weight": "0"}),
	})

	want := []string{"node-a", "node-b", "node-a", "node-b"}
	got := make([]string, 0, len(want))
	for range want {
		n, err := rr.Select("k")
		if err != nil {
			t.Fatalf("Select returned error: %v", err)
		}
		got = append(got, n.ID())
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sequence mismatch: got %v want %v", got, want)
	}
}

func TestRoundRobinSelector_UpdateResetsState(t *testing.T) {
	rr := New()
	rr.Update([]selector.Node{
		selector.NewNode("node-a", "service1", "v1.0.0", map[string]string{"weight": "0"}),
		selector.NewNode("node-b", "service1", "v1.0.0", map[string]string{"weight": "0"}),
	})

	for i := 0; i < 3; i++ {
		if _, err := rr.Select("k"); err != nil {
			t.Fatalf("Select returned error: %v", err)
		}
	}

	rr.Update([]selector.Node{
		selector.NewNode("node-c", "service1", "v1.0.0", map[string]string{"weight": "0"}),
		selector.NewNode("node-d", "service1", "v1.0.0", map[string]string{"weight": "0"}),
	})

	want := []string{"node-c", "node-d", "node-c", "node-d"}
	got := make([]string, 0, len(want))
	for range want {
		n, err := rr.Select("k")
		if err != nil {
			t.Fatalf("Select returned error: %v", err)
		}
		got = append(got, n.ID())
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sequence mismatch after Update: got %v want %v", got, want)
	}
}

func TestRoundRobinSelector_Select_WithFilter(t *testing.T) {
	rr := New()
	rr.Update([]selector.Node{
		selector.NewNode("node-a", "service1", "v1.0.0", map[string]string{"weight": "0"}),
		selector.NewNode("node-b", "service1", "v1.0.0", map[string]string{"weight": "0"}),
		selector.NewNode("node-c", "service1", "v1.0.0", map[string]string{"weight": "0"}),
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

	got, err := rr.Select("k", onlyB)
	if err != nil {
		t.Fatalf("Select returned error: %v", err)
	}
	if got.ID() != "node-b" {
		t.Fatalf("expected node-b, got %v", got.ID())
	}
}
