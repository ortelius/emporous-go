package collection

import (
	"github.com/uor-framework/client/model"
)

var _ model.NodeBuilder = &collectionBuilder{}

type collectionBuilder struct {
	nodes []model.Node
	edges []model.Edge
}

// NewBuilder returns a builder for collection nodes.
func NewBuilder(nodes []model.Node, edges []model.Edge) model.NodeBuilder {
	return &collectionBuilder{
		nodes: nodes,
		edges: edges,
	}
}

// Build completes any required actions for assembly
// before return the final immutable collection.
// At node build time create and attach the iterator.
func (b *collectionBuilder) Build(id string) (model.Node, error) {
	c := NewCollection(id)
	for _, node := range b.nodes {
		if err := c.AddNode(node); err != nil {
			return nil, err
		}
	}
	for _, edge := range b.edges {
		if err := c.AddEdge(edge); err != nil {
			return nil, err
		}
	}
	itr := NewByAttributesIterator(c.Nodes())
	c.ByAttributesIterator = itr
	return c, nil
}
