package graph

type NodeVisitor interface {
	VisitBuilderNode(*BuildNode)
}
