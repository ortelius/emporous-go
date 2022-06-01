package graph

type Node interface {
	// ID returns the unique ID of the Node
	ID() string
	// Accept will allow a visitor to access node data
	Accept(NodeVisitor)
}
