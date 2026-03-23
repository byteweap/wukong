package roundrobin

import (
	"errors"
	"reflect"
	"testing"

	"github.com/byteweap/wukong/component/selector"
)

type testNode struct {
	id     string
	weight float64
	meta   map[string]any
}

func (n testNode) ID() string           { return n.id }
func (n testNode) App() string          { return "" }
func (n testNode) Weight() float64      { return n.weight }
func (n testNode) Scheme() string       { return "" }
func (n testNode) Version() string      { return "" }
func (n testNode) Meta() map[string]any { return n.meta }

func TestRoundRobinSelector_Select_EmptyNodes(t *testing.T) {
	rr := NewRoundRobinSelector()
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
	rr := NewRoundRobinSelector()
	nodes := []selector.Node{
		testNode{id: "node-a", weight: 1, meta: map[string]any{"zone": "1"}},
		testNode{id: "node-b", weight: 2, meta: map[string]any{"zone": "2"}},
	}

	rr.Update(nodes)
	got := rr.Nodes()
	if len(got) != len(nodes) || got[0].ID() != "node-a" || got[1].ID() != "node-b" {
		t.Fatalf("Nodes() mismatch: got %#v want %#v", got, nodes)
	}
}

func TestRoundRobinSelector_Select_Sequence(t *testing.T) {
	rr := NewRoundRobinSelector()
	rr.Update([]selector.Node{
		testNode{id: "node-a"},
		testNode{id: "node-b"},
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
	rr := NewRoundRobinSelector()
	rr.Update([]selector.Node{
		testNode{id: "node-a"},
		testNode{id: "node-b"},
	})

	for i := 0; i < 3; i++ {
		if _, err := rr.Select("k"); err != nil {
			t.Fatalf("Select returned error: %v", err)
		}
	}

	rr.Update([]selector.Node{
		testNode{id: "node-c"},
		testNode{id: "node-d"},
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
	rr := NewRoundRobinSelector()
	rr.Update([]selector.Node{
		testNode{id: "node-a"},
		testNode{id: "node-b"},
		testNode{id: "node-c"},
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
