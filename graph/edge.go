package graph

type Edge struct {
	From Node
	To   Node
}

func NewEdge(from, to Node) Edge {
	return Edge{From: from, To: to}
}
