package graph

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"text/template"
)

// Node defines a single unit containing build information about a file.
type Node struct {
	// Unique node name
	Name string
	// Nodes will describe nodes connected to this one
	Nodes map[string]*Node
	// Builder specific fields
	Template template.Template
	Links    map[string]interface{}
}

// NewNode create a empty Node.
func NewNode(name string) *Node {
	return &Node{
		Name:  name,
		Nodes: map[string]*Node{},
		Links: map[string]interface{}{},
	}
}

// Graph defines a collection of Nodes.
type Graph struct {
	// Nodes describes all nodes contained in the graph
	Nodes map[string]*Node
}

// NewGraph creates an empty Graph.
func NewGraph() *Graph {
	return &Graph{
		Nodes: map[string]*Node{},
	}
}

// AddNode adds a new node to the graph.
func (g *Graph) AddNode(name string) {
	n := NewNode(name)
	g.Nodes[name] = n
}

// AddNodeTemplate adds a template to the node at the specified key in the graph.
func (g *Graph) AddNodeTemplate(key string, t template.Template) error {
	n, found := g.Nodes[key]
	if !found {
		return fmt.Errorf("node %v not found in graph", key)
	}
	n.Template = t
	g.Nodes[key] = n
	return nil
}

// AddNodeLinkInformation adds link data to the node at the specified key in the graph.
func (g *Graph) AddNodeLinkInformation(key string, links map[string]interface{}) error {
	n, found := g.Nodes[key]
	if !found {
		return fmt.Errorf("node %v not found in graph", key)
	}
	n.Links = links
	g.Nodes[key] = n
	return nil
}

// AddEdge adds an edge between two nodes in the graph
func (g *Graph) AddEdge(origin, destination string) error {
	n1 := g.Nodes[origin]
	n2 := g.Nodes[destination]

	// return an error if one of the nodes doesn't exist
	if n1 == nil || n2 == nil {
		return errors.New("not all nodes exist")
	}

	// do nothing if the node are already connected
	if _, ok := n1.Nodes[n2.Name]; ok {
		return nil
	}

	n1.Nodes[n2.Name] = n2

	// Add the nodes to the graph's node map
	g.Nodes[n1.Name] = n1
	g.Nodes[n2.Name] = n2

	return nil
}

// Root calculates to root node of the graph.
// This is calculated based on existing child nodes.
// This expects only one root node to be found.
func (g *Graph) Root() (*Node, error) {
	// FIXME(jpowe432): Optimize or redesign the chain

	childNodes := map[string]int{}
	for _, n := range g.Nodes {
		for _, ch := range n.Nodes {
			childNodes[ch.Name]++
		}
	}
	var roots []*Node
	for _, n := range g.Nodes {
		if _, found := childNodes[n.Name]; !found {
			roots = append(roots, n)
		}
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("no root found in graph")
	}
	if len(roots) > 1 {
		var rootNames []string
		for _, root := range roots {
			rootNames = append(rootNames, root.Name)
		}
		sort.Strings(rootNames)
		return nil, fmt.Errorf("multiple roots found in graph: %s", strings.Join(rootNames, ", "))
	}
	return roots[0], nil
}
