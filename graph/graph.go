package graph

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Graph represent a dataset and the relationship between objects.
// This can consist of one or more OCI artifact(s).
type Graph struct {
	// Nodes describes all nodes contained in the graph
	Nodes map[string]Node
	// From describes all edges with the
	// origin node as the map key.
	From map[string]map[string]Edge
	// To describes all edges with the
	// destination node as the map key
	To map[string]map[string]Edge
}

// NewGraph creates an empty Graph.
func NewGraph() *Graph {
	return &Graph{
		Nodes: map[string]Node{},
		From:  map[string]map[string]Edge{},
		To:    map[string]map[string]Edge{},
	}
}

// AddNode adds a new node to the graph.
func (g *Graph) AddNode(node Node) error {
	if _, exists := g.Nodes[node.ID()]; exists {
		return errors.New("node ID collision")
	}
	g.Nodes[node.ID()] = node
	return nil
}

// AddEdge adds an edge between two nodes in the graph
func (g *Graph) AddEdge(from, to string) error {
	if from == to {
		return errors.New("adding self edge")
	}

	n1 := g.Nodes[from]
	n2 := g.Nodes[to]

	// return an error if one of the nodes doesn't exist
	if n1 == nil || n2 == nil {
		return errors.New("not all nodes exist")
	}

	if g.ConnectedFrom(from, to) && g.ConnectedTo(from, to) {
		return nil
	}

	edge := NewEdge(n1, n2)

	g.setEdgeFrom(edge)
	g.setEdgeTo(edge)

	return nil
}

func (g *Graph) setEdgeFrom(edge Edge) {
	from, ok := g.From[edge.From.ID()]
	if !ok {
		g.From[edge.From.ID()] = map[string]Edge{edge.To.ID(): edge}
		return
	}
	from[edge.To.ID()] = edge
	g.From[edge.From.ID()] = from
}

func (g *Graph) setEdgeTo(edge Edge) {
	to, ok := g.To[edge.To.ID()]
	if !ok {
		g.To[edge.To.ID()] = map[string]Edge{edge.From.ID(): edge}
		return
	}
	to[edge.From.ID()] = edge
	g.To[edge.To.ID()] = to
}

// Node returns the node based on the ID if the node exists.
func (g *Graph) Node(id string) Node {
	node, ok := g.Nodes[id]
	if !ok {
		return nil
	}
	return node
}

// Edge returns the edge from the origin to the destination node if such an edge exists.
// The node from must be directly reachable from the node to as defined by the From method.
func (g *Graph) Edge(from, to string) Edge {
	edge, ok := g.From[from][to]
	if !ok {
		return Edge{}
	}
	return edge
}

// ConnectedFrom returns whether there is an edge
// from the origin to the destination Node.
func (g *Graph) ConnectedFrom(from, to string) bool {
	_, ok := g.From[to][from]
	return ok
}

// ConnectedTo returns whether there is an edge
// from the destination to the origin Node.
func (g *Graph) ConnectedTo(from, to string) bool {
	_, ok := g.To[to][from]
	return ok
}

// NodesFrom returns a list of Nodes connected
// to the node with the id.
func (g *Graph) NodesFrom(id string) []Node {
	var connectedNodes []Node
	nodes, ok := g.From[id]
	if !ok {
		return nil
	}
	for id := range nodes {
		connectedNodes = append(connectedNodes, g.Nodes[id])
	}

	return connectedNodes
}

// NodesTo returns a list of Nodes connected
// to the node with the id.
func (g *Graph) NodesTo(id string) []Node {
	var connectedNodes []Node
	nodes, ok := g.To[id]
	if !ok {
		return nil
	}
	for id := range nodes {
		connectedNodes = append(connectedNodes, g.Nodes[id])
	}

	return connectedNodes
}

// Root calculates to root node of the graph.
// This is calculated base on existing child nodes.
// This expects only one root node to be found.
func (g *Graph) Root() (Node, error) {
	childNodes := map[string]int{}
	for _, n := range g.Nodes {
		for _, ch := range g.NodesFrom(n.ID()) {
			childNodes[ch.ID()]++
		}
	}
	var roots []Node
	for _, n := range g.Nodes {
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
			rootNames = append(rootNames, root.ID())
		}
		sort.Strings(rootNames)
		return nil, fmt.Errorf("multiple roots found in graph: %s", strings.Join(rootNames, ", "))
	}
	return roots[0], nil
}
