package collection

import (
	"errors"

	"github.com/uor-framework/uor-client-go/model"
)

// ErrNodesNotExist is an error that is thrown if an edge
// is added without existing nodes.
var ErrNodesNotExist = errors.New("not all nodes exist")

// AddNode adds a new node to the graph.
func (c *Collection) AddNode(node model.Node) error {
	if _, exists := c.nodes[node.ID()]; exists {
		return errors.New("node ID collision")
	}
	c.nodes[node.ID()] = node
	return nil
}

// UpdateNode adds a new node or updates the existing node, if applicable.
func (c *Collection) UpdateNode(node model.Node) error {
	c.nodes[node.ID()] = node
	return nil
}

// AddEdge adds an edge between two nodes in the graph
func (c *Collection) AddEdge(edge model.Edge) error {
	from := edge.From().ID()
	to := edge.To().ID()

	if from == to {
		return errors.New("adding self edge")
	}

	n1 := c.nodes[from]
	n2 := c.nodes[to]

	// return an error if one of the nodes doesn't exist
	if n1 == nil || n2 == nil {
		return ErrNodesNotExist
	}

	if c.HasEdgeFromTo(from, to) {
		return nil
	}

	c.setEdgeFrom(edge)
	c.setEdgeTo(edge)

	return nil
}

// SubCollection returns a sub-collection with only the nodes that satisfy the matcher.
func (c *Collection) SubCollection(matcher model.Matcher) (Collection, error) {
	if matcher == nil {
		return *c, nil
	}
	out := New(c.ID())
	out.Location = c.Address()
	for _, node := range c.Nodes() {
		match, err := matcher.Matches(node)
		if err != nil {
			return *out, err
		}
		if match {
			if err := out.AddNode(node); err != nil {
				return Collection{}, err
			}
		}
	}

	for _, edge := range c.Edges() {
		err := out.AddEdge(edge)
		if err != nil && !errors.Is(err, ErrNodesNotExist) {
			return *out, err
		}
	}
	return *out, nil
}

func (c *Collection) setEdgeFrom(edge model.Edge) {
	from, ok := c.from[edge.From().ID()]
	if !ok {
		c.from[edge.From().ID()] = map[string]model.Edge{edge.To().ID(): edge}
		return
	}
	from[edge.To().ID()] = edge
	c.from[edge.From().ID()] = from
}

func (c *Collection) setEdgeTo(edge model.Edge) {
	to, ok := c.to[edge.To().ID()]
	if !ok {
		c.to[edge.To().ID()] = map[string]model.Edge{edge.From().ID(): edge}
		return
	}
	to[edge.From().ID()] = edge
	c.to[edge.To().ID()] = to
}
