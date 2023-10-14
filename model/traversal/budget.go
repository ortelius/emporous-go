package traversal

import (
	"fmt"

	"github.com/emporous/emporous-go/model"
)

// Similar to the budget in the go-ipld-prime library.
// Reference: https://github.com/ipld/go-ipld-prime/blob/ab0f17bec1e700e4c76a6bbc28e7260cea7c035d/traversal/fns.go#L114

// Budget tracks budgeted operations during graph traversal.
type Budget struct {
	// Maximum numbers of nodes to visit in a single traversal operator before stopping.
	NodeBudget int64
}

// ErrBudgetExceeded is an error that described the event where
// the maximum amount of nodes have been visited with no match.
type ErrBudgetExceeded struct {
	Node model.Node
}

func (e *ErrBudgetExceeded) Error() string {
	return fmt.Sprintf("traversal budget exceeded: node budget for reached zero while on node %v", e.Node.Address())
}
