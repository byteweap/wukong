package selector

import (
	"testing"
)

func TestErrNoAvailableNode(t *testing.T) {
	if ErrNoAvailableNode.Error() != "no available node" {
		t.Errorf("unexpected error message: %s", ErrNoAvailableNode.Error())
	}
}

func TestNewNode(t *testing.T) {
	meta := map[string]string{"weight": "10"}
	node := NewNode("id1", "service1", "v1", meta)

	if node.ID() != "id1" {
		t.Errorf("expected id1, got %s", node.ID())
	}
	if node.Service() != "service1" {
		t.Errorf("expected service1, got %s", node.Service())
	}
	if node.Version() != "v1" {
		t.Errorf("expected v1, got %s", node.Version())
	}
	if node.Weight() != 10.0 {
		t.Errorf("expected 10.0, got %f", node.Weight())
	}
	if node.Meta()["weight"] != "10" {
		t.Errorf("expected 10, got %s", node.Meta()["weight"])
	}
}

func TestNode_Weight(t *testing.T) {
	t.Run("with weight", func(t *testing.T) {
		node := NewNode("id1", "service1", "v1", map[string]string{"weight": "5.5"})
		if node.Weight() != 5.5 {
			t.Errorf("expected 5.5, got %f", node.Weight())
		}
	})

	t.Run("without weight", func(t *testing.T) {
		node := NewNode("id1", "service1", "v1", nil)
		if node.Weight() != 0 {
			t.Errorf("expected 0, got %f", node.Weight())
		}
	})

	t.Run("invalid weight", func(t *testing.T) {
		node := NewNode("id1", "service1", "v1", map[string]string{"weight": "invalid"})
		if node.Weight() != 0 {
			t.Errorf("expected 0, got %f", node.Weight())
		}
	})
}

func TestVersion_Filter(t *testing.T) {
	nodes := []Node{
		NewNode("1", "s", "v1", nil),
		NewNode("2", "s", "v2", nil),
		NewNode("3", "s", "v1", nil),
	}

	filter := Version("v1")
	result := filter(nodes)

	if len(result) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(result))
	}
	for _, n := range result {
		if n.Version() != "v1" {
			t.Errorf("expected v1, got %s", n.Version())
		}
	}
}

func TestVersion_Filter_NoMatch(t *testing.T) {
	nodes := []Node{
		NewNode("1", "s", "v1", nil),
		NewNode("2", "s", "v2", nil),
	}

	filter := Version("v3")
	result := filter(nodes)

	if len(result) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(result))
	}
}

func TestVersion_Filter_EmptyNodes(t *testing.T) {
	filter := Version("v1")
	result := filter([]Node{})

	if len(result) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(result))
	}
}
