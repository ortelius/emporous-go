package collection

import (
	"fmt"
	"sort"
	"strings"

	"github.com/uor-framework/uor-client-go/model"
)

var (
	_ model.Node     = &Collection{}
	_ model.Rooted   = &Collection{}
	_ model.Iterator = &Collection{}
)

// Collection is implementation of a model Node represent one OCI artifact.
type Collection struct {
	// unique ID for the collection
	id string
	// nodes describes all nodes contained in the graph
	nodes map[string]model.Node
	// from describes all edges with the
	// origin node as the map key.
	from map[string]map[string]model.Edge
	// to describes all edges with the
	// destination node as the map key
	to map[string]map[string]model.Edge
	// Location of the collection (local or remote)
	Location string
	// Iterator for the collection node
	*ByAttributesIterator
}

// New creates an empty Collection with the specified ID.
func New(id string) *Collection {
	return &Collection{
		id:                   id,
		nodes:                map[string]model.Node{},
		from:                 map[string]map[string]model.Edge{},
		to:                   map[string]map[string]model.Edge{},
		ByAttributesIterator: NewByAttributesIterator(nil),
	}
}

// ID return the unique id of the collection.
func (c *Collection) ID() string {
	return c.id
}

// Address returns collection location.
func (c *Collection) Address() string {
	return c.Location
}

// Attributes returns a collection of all the
// attributes contained within the collection nodes.
// Because each parent node should inherit the attributes, all
// the attached child nodes, the root node will contain attributes
// for the entire collection. If no root node exists, nil is returned.
func (c *Collection) Attributes() model.AttributeSet {
	root, err := c.Root()
	if err == nil {
		return root.Attributes()
	}
	return nil
}

// NodeByID returns the node based on the ID if the node exists.
func (c *Collection) NodeByID(id string) model.Node {
	node, ok := c.nodes[id]
	if !ok {
		return nil
	}
	return node
}

// Nodes returns a slice containing
// all nodes in the graph.
func (c *Collection) Nodes() []model.Node {
	var nodes []model.Node
	for _, node := range c.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// Edges returns a slice containing
// all nodes in the graph.
func (c *Collection) Edges() []model.Edge {
	var edges []model.Edge
	for _, to := range c.from {
		for _, edge := range to {
			edges = append(edges, edge)
		}
	}
	return edges
}

// Edge returns the edge from the origin to the destination node if such an edge exists.
// The node from must be directly reachable from the node to as defined by the From method.
func (c *Collection) Edge(from, to string) model.Edge {
	edge, ok := c.from[from][to]
	if !ok {
		return nil
	}
	return edge
}

// HasEdgeFromTo returns whether there is an edge
// from the origin to the destination Node.
func (c *Collection) HasEdgeFromTo(from, to string) bool {
	_, ok := c.from[to][from]
	return ok
}

// From returns a list of Nodes connected
// to the node with the id.
func (c *Collection) From(id string) []model.Node {
	var connectedNodes []model.Node
	nodes, ok := c.from[id]
	if !ok {
		return nil
	}
	for id := range nodes {
		connectedNodes = append(connectedNodes, c.nodes[id])
	}
	return connectedNodes
}

// To returns a list of Nodes connected
// to the node with the id.
func (c *Collection) To(id string) []model.Node {
	var connectedNodes []model.Node
	nodes, ok := c.to[id]
	if !ok {
		return nil
	}
	for id := range nodes {
		connectedNodes = append(connectedNodes, c.nodes[id])
	}

	return connectedNodes
}

// Root calculates to root node of the graph.
// This is calculated base on existing child nodes.
// This expects only one root node to be found.
func (c *Collection) Root() (model.Node, error) {
	childNodes := map[string]int{}
	for _, n := range c.nodes {
		for _, ch := range c.From(n.ID()) {
			childNodes[ch.ID()]++
		}
	}
	var roots []model.Node
	for _, n := range c.nodes {
		if _, found := childNodes[n.ID()]; !found {
			roots = append(roots, n)
		}
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("no root found in graph")
	}
	if len(roots) > 1 {
		var rootNames []string
		for _, root := range roots {
			rootNames = append(rootNames, root.Address())
		}
		sort.Strings(rootNames)
		return nil, fmt.Errorf("multiple roots found in graph: %s", strings.Join(rootNames, ", "))
	}
	return roots[0], nil
}
