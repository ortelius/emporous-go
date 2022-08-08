package descriptor

import (
	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// Node defines a single unit containing information about a UOR dataset node.
type Node struct {
	id         string
	descriptor ocispec.Descriptor
	attributes model.AttributeSet
	Location   string
}

var _ model.Node = &Node{}

// NewNode create an empty Descriptor Node.
func NewNode(id string, descriptor ocispec.Descriptor) *Node {
	attr := AnnotationsToAttributes(descriptor.Annotations)
	return &Node{
		id:         id,
		attributes: attr,
		descriptor: descriptor,
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
func (n *Node) Attributes() model.AttributeSet {
	return n.attributes
}

// Descriptor returns the underlying descriptor object.
func (n *Node) Descriptor() ocispec.Descriptor {
	return n.descriptor
}

// AnnotationsToAttributes converts annotations from a descriptors
// to an Attribute type. Any value that is not valid JSON will be skipped.
func AnnotationsToAttributes(annotations map[string]string) model.AttributeSet {
	attr := attributes.Attributes{}
	for key, value := range annotations {
		attr[key] = attributes.NewString(key, value)
	}
	return attr
}
