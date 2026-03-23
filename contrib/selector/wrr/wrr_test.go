package wrr

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
		testNode{id: "node-a", weight: 1, meta: map[string]any{"zone": "1"}},
		testNode{id: "node-b", weight: 2, meta: map[string]any{"zone": "2"}},
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
		testNode{id: "node-a", weight: 2},
		testNode{id: "node-b", weight: 1},
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
		testNode{id: "node-a", weight: 2},
		testNode{id: "node-b", weight: 1},
	})

	for i := 0; i < 3; i++ {
		if _, err := ws.Select("ignored"); err != nil {
			t.Fatalf("Select returned error: %v", err)
		}
	}

	ws.Update([]selector.Node{
		testNode{id: "node-c", weight: 1},
		testNode{id: "node-d", weight: 1},
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
		testNode{id: "node-a", weight: -2},
		testNode{id: "node-b", weight: 1},
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
		testNode{id: "node-a", weight: 1},
		testNode{id: "node-b", weight: 1},
		testNode{id: "node-c", weight: 1},
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
