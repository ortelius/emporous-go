package collection

import (
	"sort"

	"github.com/uor-framework/client/model"
)

var (
	_ model.Iterator = &InOrderIterator{}
	_ model.Iterator = &ByAttributesIterator{}
)

// InOrderIterator implements the model.Iterator interface and traverse the nodes
// in the order provided.
type InOrderIterator struct {
	idx   int
	nodes []model.Node
}

// NewInOrderIterator returns a OrderedNodes initialized with the provided nodes.
func NewInOrderIterator(nodes []model.Node) *InOrderIterator {
	return &InOrderIterator{idx: -1, nodes: nodes}
}

// Len returns the remaining number of nodes to be iterated over.
func (n *InOrderIterator) Len() int {
	if n.idx >= len(n.nodes) {
		return 0
	}
	return len(n.nodes[n.idx+1:])
}

// Next returns whether the next call of Node will return a valid node.
func (n *InOrderIterator) Next() bool {
	if uint(n.idx)+1 < uint(len(n.nodes)) {
		n.idx++
		return true
	}
	n.idx = len(n.nodes)
	return false
}

// Node returns the current node of the iterator. Next must have been
// called prior to a call to Node.
func (n *InOrderIterator) Node() model.Node {
	if n.idx >= len(n.nodes) || n.idx < 0 {
		return nil
	}
	return n.nodes[n.idx]
}

// Reset returns the iterator to its initial state.
func (n *InOrderIterator) Reset() {
	n.idx = -1
}

// Error returns found errors during iteration.
func (n *InOrderIterator) Error() error {
	return nil
}

// LazyOrderedNodesByAttribute implements the model.Iterator interface and traverse the nodes
// in from smallest to largest attribute list.
type ByAttributesIterator struct {
	iter       InOrderIterator
	attributes ByAttributeSet
}

// NewLazyOrderedNodesByAttribute returns a LazyOrderedNodesAttribute initialized with the
// provided nodes.
func NewByAttributesIterator(nodes []model.Node) *ByAttributesIterator {
	return &ByAttributesIterator{attributes: nodes}
}

// Len returns the remaining number of nodes to be iterated over.
func (n *ByAttributesIterator) Len() int {
	if n.iter.nodes == nil {
		return len(n.attributes)
	}
	return n.iter.Len()
}

// Next returns whether the next call of Node will return a valid node.
func (n *ByAttributesIterator) Next() bool {
	if n.iter.nodes == nil {
		n.fillSlice()
	}
	return n.iter.Next()
}

// Node returns the current node of the iterator. Next must have been
// called prior to a call to Node.
func (n *ByAttributesIterator) Node() model.Node {
	return n.iter.Node()
}

// Reset returns the iterator to its initial state.
func (n *ByAttributesIterator) Reset() {
	n.iter.Reset()
}

// Error returns found errors during iteration.
func (n *ByAttributesIterator) Error() error {
	return n.iter.Error()
}

func (n *ByAttributesIterator) fillSlice() {
	sort.Sort(n.attributes)
	n.iter = InOrderIterator{
		idx:   -1,
		nodes: n.attributes,
	}
	n.attributes = nil
}

// ByAttributeSet is a slice of Node sorted by attribute set length.
type ByAttributeSet []model.Node

func (a ByAttributeSet) Len() int           { return len(a) }
func (a ByAttributeSet) Less(i, j int) bool { return a[i].Attributes().Len() < a[j].Attributes().Len() }
func (a ByAttributeSet) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
