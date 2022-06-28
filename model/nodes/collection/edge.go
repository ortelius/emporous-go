package collection

import "github.com/uor-framework/client/model"

// Edge defines a relationship
// between two Nodes.
type Edge struct {
	F model.Node
	T model.Node
}

// NewEdge returns a new Edge instance.
func NewEdge(from, to model.Node) Edge {
	return Edge{F: from, T: to}
}

// To return the destination node
// for the edge.
func (e Edge) To() model.Node {
	return e.T
}

// From returns the origin node
// for the edge
func (e Edge) From() model.Node {
	return e.F
}
