package model

import "encoding/json"

// DirectedGraph defines methods for interacting with groups of
// nodes and edges in a directed graph structure. This graph may
// or may not contain cycles.
type DirectedGraph interface {
	// Nodes lists all nodes contained within the
	// graph.
	Nodes() []Node
	// Edges lists all edges contained within the
	// graph.
	Edges() []Edge
	// NodeByID returns a node from the given
	// id, if existing.
	NodeByID(string) Node
	// Edge creates or returns an Edge from the given
	// node ids.
	Edge(string, string) Edge
	// HasEdgeFromTo returns whether an edge exits between
	// two nodes.
	HasEdgeFromTo(string, string) bool
	// From lists all child nodes for a given node
	// id, if existing.
	From(string) []Node
	// To lists all parent nodes for a given node
	// id, if existing.
	To(string) []Node
}

// Rooted defines methods to locate the root of the
// data set.
type Rooted interface {
	Root() (Node, error)
}

// Node defines read-only methods implemented by different
// node types.
type Node interface {
	// ID is a unique value assigned to the node.
	ID() string
	// Address is the location where the data is stored
	Address() string
	// Attributes defines the attributes associated
	// with the node data
	Attributes() AttributeSet
}

// NodeBuilder defines methods to build new nodes.
type NodeBuilder interface {
	// Build create a new immutable node.
	Build(string) (Node, error)
}

// Edge defines methods for node relationship
// information.
// This may eventually include weight information to
// represent nodes at addresses that are a longer
// distance from the source.
type Edge interface {
	// To is the destination node.
	To() Node
	// From is the origin node.
	From() Node
}

// Iterator defines method for traversing node data in
// a specified order.
type Iterator interface {
	// Next returns true if there is more data to iterate
	// and will increment.
	Next() bool
	// Node will return the node in the current position.
	Node() Node
	// Reset will start the iterator from the beginning
	Reset()
	// Error will return all accumulated errors during iteration.
	Error() error
}

// Matcher defines methods used for node searching.
type Matcher interface {
	// Matches evaluates the current node against the criteria.
	Matches(node Node) bool
}

// AttributeSet defines methods for manipulating attribute sets.
// Nodes have set of attributes that allow them to self-describe and
// describe connected nodes.
type AttributeSet interface {
	// Exists returns whether a key, value with type pair exists
	Exists(string, Kind, interface{}) bool
	// Find returns all values associated with a specified key
	Find(string) Attribute
	// AsJSON returns a json representation of the Attribute set.
	AsJSON() json.RawMessage
	// List will list all key,value pairs for the attributes in a
	// consumable format.
	List() map[string]Attribute
	// Len returns the attribute set length
	Len() int
}

// Attribute defines methods of an attribute object.
type Attribute interface {
	// Key is the value of the attribute identifier. This must always be a string.
	Key() string
	// Kind represent the value type.
	Kind() Kind
	// IsNull will return true if the attribute type is null.
	IsNull() bool
	// AsBool will return the attribute values as a boolean.
	AsBool() (bool, error)
	// AsNumber will return the attribute value as a float.
	AsNumber() (float64, error)
	// AsString will return the attribute value as a string.
	AsString() (string, error)
	// AsAny returns the value of the attribute with no type checking.
	AsAny() interface{}
}

// Kind represents the kind of Attributes.
type Kind int

const (
	KindInvalid Kind = iota
	KindNull
	KindBool
	KindNumber
	KindString
)

// String prints a string representation of the attribute kind.
func (k Kind) String() string {
	switch k {
	case KindInvalid:
		return "INVALID"
	case KindNull:
		return "null"
	case KindBool:
		return "bool"
	case KindNumber:
		return "number"
	case KindString:
		return "string"
	default:
		panic("invalid kind")
	}
}
