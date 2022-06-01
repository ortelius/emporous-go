package graph

// BuildNode defines a single unit containing build information about a file for building.
type BuildNode struct {
	// Unique node name
	Name string
}

var _ Node = &BuildNode{}

// NewNode create a empty Node.
func NewBuildNode(name string) *BuildNode {
	return &BuildNode{
		Name: name,
	}
}

func (n *BuildNode) ID() string {
	return n.Name
}

func (n *BuildNode) Accept(v NodeVisitor) {
	v.VisitBuilderNode(n)
}
