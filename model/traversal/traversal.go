package traversal

import (
	"context"
	"errors"

	"github.com/uor-framework/uor-client-go/model"
)

// ErrSkip allows a node to be intentionally skipped.
var ErrSkip = errors.New("skip")

// Tracker contains information stored during graph traversal
// such as the node path, budgeting information, and visited nodes.
type Tracker struct {
	// Path defines a series of steps during traversal over nodes.
	Path
	// Budget tracks traversal maximums such as maximum visited nodes.
	budget *Budget
}

// NewTracker returns a new Tracker instance.
func NewTracker(root model.Node, budget *Budget) Tracker {
	t := Tracker{
		budget: budget,
		Path:   NewPath(root),
	}
	return t
}

// Walk the nodes of a graph and call the handler for each. If the handler
// decodes the child nodes for each parent node, they are visited as well.
func Walk(ctx context.Context, handler Handler, node model.Node) error {
	tracker := NewTracker(node, nil)
	return tracker.Walk(ctx, handler, node)
}

// Walk the nodes of a graph  and call the handler for each. If the handler
// decodes the child nodes for each parent node, they will be visited as well.
// The node budget and path traversal steps are tracked with the Tracker.
// This function is based on github.com/containerd/containerd/images.Walk.
func (t Tracker) Walk(ctx context.Context, handler Handler, nodes ...model.Node) error {
	for _, node := range nodes {

		if t.budget != nil {
			if t.budget.NodeBudget <= 0 {
				return &ErrBudgetExceeded{Node: node}
			}
			t.budget.NodeBudget--
		}

		children, err := handler.Handle(ctx, t, node)
		if err != nil {
			if errors.Is(err, ErrSkip) {
				continue // don't traverse the children.
			}
			return err
		}

		tNext := t
		for _, child := range children {
			tNext.Path = t.Path.Add(node, child)
		}

		if len(children) > 0 {
			if err := tNext.Walk(ctx, handler, children...); err != nil {
				return err
			}
		}
	}
	return nil
}
