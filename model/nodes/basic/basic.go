package basic

import (
	"github.com/uor-framework/uor-client-go/model"
)

// Node defines a single unit containing information about a UOR dataset node.
type Node struct {
	id         string
	attributes model.Attributes
	Location   string
}

var _ model.Node = &Node{}

// NewNode create an empty Basic Node.
func NewNode(id string, attributes model.Attributes) *Node {
	return &Node{
		id:         id,
		attributes: attributes,
	}
}

// ID returns the unique identifier for a  basic Node.
func (n *Node) ID() string {
	return n.id
}

// Address returns the set location for basic Node
// data.
func (n *Node) Address() string {
	return n.Location
}

// Attributes represents a collection of data defining the node.
func (n *Node) Attributes() model.Attributes {
	return n.attributes
}
