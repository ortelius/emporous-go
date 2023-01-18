package testutils

import (
	"github.com/emporous/emporous-go/model"
)

var (
	_ model.Node = &FakeNode{}
	_ model.Node = &FakeIterableNode{}
)

// FakeNode implements the model.Node interface for testing.
type FakeNode struct {
	// Node ID
	I string
	// Node Attributes
	A model.AttributeSet
}

func (m *FakeNode) ID() string {
	return m.I
}

func (m *FakeNode) Address() string {
	return "address"
}

func (m *FakeNode) Attributes() model.AttributeSet {
	return m.A
}

// FakeIterableNode implements the model.Node and model.Iterator interface for testing.
type FakeIterableNode struct {
	// Node ID
	I string
	// Iterator Index
	Index int
	// Node Attributes
	A model.AttributeSet
	// Iterable nodes list
	Nodes []model.Node
}

func (m *FakeIterableNode) ID() string {
	return m.I
}

func (m *FakeIterableNode) Address() string {
	return "address"
}

func (m *FakeIterableNode) Attributes() model.AttributeSet {
	return m.A
}

func (m *FakeIterableNode) Len() int {
	if m.Index >= len(m.Nodes) {
		return 0
	}
	return len(m.Nodes[m.Index+1:])
}

func (m *FakeIterableNode) Next() bool {
	if uint(m.Index)+1 < uint(len(m.Nodes)) {
		m.Index++
		return true
	}
	m.Index = len(m.Nodes)
	return false
}

func (m *FakeIterableNode) Node() model.Node {
	if m.Index >= len(m.Nodes) || m.Index < 0 {
		return nil
	}
	return m.Nodes[m.Index]
}

func (m *FakeIterableNode) Reset() {
	m.Index = -1
}

func (m *FakeIterableNode) Error() error {
	return nil
}
