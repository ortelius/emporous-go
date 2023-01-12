package traversal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/emporous/emporous-go/model"
	"github.com/emporous/emporous-go/util/testutils"
)

func TestTracker_Walk(t *testing.T) {
	type spec struct {
		name           string
		t              Tracker
		root           model.Node
		graph          model.DirectedGraph
		expError       error
		expInvocations int
	}

	cases := []spec{
		{
			name: "Success/VisitRootNode",

			t: Tracker{
				budget: &Budget{
					NodeBudget: 3,
				},
				Path: NewPath(&testutils.FakeNode{I: "node1"}),
			},
			root: &testutils.FakeNode{I: "node1"},
			graph: &mockGraph{nodes: map[string][]model.Node{
				"node1": {&testutils.FakeNode{I: "node2"}}},
			},
			expInvocations: 2,
		},
		{
			name: "Success/DuplicateNodeID",
			t: Tracker{
				budget: &Budget{
					NodeBudget: 8,
				},
				Path: NewPath(&testutils.FakeNode{I: "node1"}),
			},
			root: &testutils.FakeNode{I: "node1"},
			graph: &mockGraph{nodes: map[string][]model.Node{
				"node1": {
					&testutils.FakeIterableNode{
						I:     "node2",
						Index: -1,
						Nodes: []model.Node{&testutils.FakeNode{I: "node1"}}},
				},
			},
			},
			expInvocations: 2,
		},
		{
			name: "Failure/ExceededBudget",

			t: Tracker{
				budget: &Budget{
					NodeBudget: 0,
				},
				Path: NewPath(&testutils.FakeNode{I: "node1"}),
			},
			root: &testutils.FakeNode{I: "node1"},
			graph: &mockGraph{nodes: map[string][]model.Node{
				"node1": {&testutils.FakeNode{I: "node2"}}},
			},
			expInvocations: 0,
			expError:       &ErrBudgetExceeded{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var actualInvocations int
			handler := HandlerFunc(func(ctx context.Context, tracker Tracker, node model.Node) ([]model.Node, error) {
				t.Log("Visiting " + node.ID())
				actualInvocations++
				return c.graph.From(node.ID()), nil
			})

			err := c.t.Walk(context.Background(), handler, c.root)
			if c.expError != nil {
				require.ErrorAs(t, err, &c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expInvocations, actualInvocations)
			}
		})
	}
}

// Mock tree structure
type mockGraph struct {
	nodes map[string][]model.Node
}

var _ model.DirectedGraph = &mockGraph{}

func (m *mockGraph) From(id string) []model.Node {
	return m.nodes[id]
}

func (m *mockGraph) To(_ string) []model.Node {
	return nil
}

func (m *mockGraph) Edge(_, _ string) model.Edge {
	return nil
}

func (m *mockGraph) Edges() []model.Edge {
	return nil
}

func (m *mockGraph) Nodes() []model.Node {
	return nil
}

func (m *mockGraph) HasEdgeFromTo(_, _ string) bool {
	return false
}

func (m *mockGraph) NodeByID(_ string) model.Node {
	return nil
}
