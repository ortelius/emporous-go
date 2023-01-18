package v2

import (
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/emporous/emporous-go/model"
	"github.com/emporous/emporous-go/nodes/descriptor"
)

// Node defines a single unit containing information about a Emporous dataset node.
type Node struct {
	id         string
	descriptor ocispec.Descriptor
	Properties *descriptor.Properties
	Location   string
}

var _ model.Node = &Node{}

// NewNode create a new Descriptor Node.
func NewNode(id string, desc ocispec.Descriptor) (*Node, error) {
	attr, err := descriptor.AnnotationsToAttributes(desc.Annotations)
	if err != nil {
		return nil, err
	}
	props, err := descriptor.Parse(attr)
	if err != nil {
		return nil, err
	}
	return &Node{
		id:         id,
		Properties: props,
		descriptor: desc,
	}, nil
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
func (n *Node) Attributes() model.AttributeSet {
	return n.Properties
}

// Descriptor returns the underlying descriptor object.
func (n *Node) Descriptor() ocispec.Descriptor {
	return n.descriptor
}
