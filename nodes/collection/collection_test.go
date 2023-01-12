package collection

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/util/testutils"
)

func TestCollection_Root(t *testing.T) {
	type spec struct {
		name     string
		nodes    []model.Node
		edges    []model.Edge
		expID    string
		expError string
	}

	cases := []spec{
		{
			name:  "Success/RootExists",
			nodes: []model.Node{&testutils.FakeNode{I: "node1"}, &testutils.FakeNode{I: "node2"}, &testutils.FakeNode{I: "node3"}},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node2"}, F: &testutils.FakeNode{I: "node1"}},
				&Edge{T: &testutils.FakeNode{I: "node1"}, F: &testutils.FakeNode{I: "node3"}},
			},
			expID: "node3",
		},
		{
			name:  "Failure/NotRootExists",
			nodes: []model.Node{&testutils.FakeNode{I: "node1"}, &testutils.FakeNode{I: "node2"}, &testutils.FakeNode{I: "node3"}},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node2"}, F: &testutils.FakeNode{I: "node1"}},
				&Edge{T: &testutils.FakeNode{I: "node1"}, F: &testutils.FakeNode{I: "node3"}},
				&Edge{T: &testutils.FakeNode{I: "node3"}, F: &testutils.FakeNode{I: "node2"}},
			},
			expError: "no root found in graph",
		},
		{
			name:  "Failure/MultipleRootsExist",
			nodes: []model.Node{&testutils.FakeNode{I: "node1"}, &testutils.FakeNode{I: "node2"}, &testutils.FakeNode{I: "node3"}},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node2"}, F: &testutils.FakeNode{I: "node1"}},
			},
			expError: "multiple roots found in graph: address, address",
		},
		{
			name:     "Failure/NoEdges",
			nodes:    []model.Node{&testutils.FakeNode{I: "node1"}, &testutils.FakeNode{I: "node2"}, &testutils.FakeNode{I: "node3"}},
			edges:    nil,
			expError: "multiple roots found in graph: address, address, address",
		},
		{
			name:     "Failure/NoNodes",
			nodes:    nil,
			edges:    nil,
			expError: "no root found in graph",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			collection := makeTestCollection(t, c.nodes, c.edges)
			root, err := collection.Root()
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expID, root.ID())
			}

		})
	}
}

func TestCollection_HasEdgeToFrom(t *testing.T) {
	type spec struct {
		name  string
		nodes []model.Node
		edges []model.Edge
		to    string
		from  string
		exp   bool
	}

	cases := []spec{
		{
			name:  "Success/EdgeExists",
			nodes: []model.Node{&testutils.FakeNode{I: "node1"}, &testutils.FakeNode{I: "node2"}, &testutils.FakeNode{I: "node3"}},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node2"}, F: &testutils.FakeNode{I: "node1"}},
				&Edge{T: &testutils.FakeNode{I: "node1"}, F: &testutils.FakeNode{I: "node3"}},
			},
			to:   "node1",
			from: "node2",
			exp:  true,
		},
		{
			name:  "Success/NoEdgeExists",
			nodes: []model.Node{&testutils.FakeNode{I: "node1"}, &testutils.FakeNode{I: "node2"}, &testutils.FakeNode{I: "node3"}},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node2"}, F: &testutils.FakeNode{I: "node1"}},
				&Edge{T: &testutils.FakeNode{I: "node1"}, F: &testutils.FakeNode{I: "node3"}},
			},
			to:   "node2",
			from: "node3",
			exp:  false,
		},
		{
			name:  "Success/EdgeExitsReverse",
			nodes: []model.Node{&testutils.FakeNode{I: "node1"}, &testutils.FakeNode{I: "node2"}, &testutils.FakeNode{I: "node3"}},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node2"}, F: &testutils.FakeNode{I: "node1"}},
				&Edge{T: &testutils.FakeNode{I: "node1"}, F: &testutils.FakeNode{I: "node3"}},
				&Edge{T: &testutils.FakeNode{I: "node3"}, F: &testutils.FakeNode{I: "node2"}},
			},
			to:   "node2",
			from: "node1",
			exp:  false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			collection := makeTestCollection(t, c.nodes, c.edges)
			actual := collection.HasEdgeFromTo(c.from, c.to)
			require.Equal(t, c.exp, actual)
		})
	}
}

func TestCollection_To(t *testing.T) {
	type spec struct {
		name   string
		nodes  []model.Node
		edges  []model.Edge
		input  string
		expIDs []string
	}

	cases := []spec{
		{
			name:  "Success/OneNodeFound",
			nodes: []model.Node{&testutils.FakeNode{I: "node1"}, &testutils.FakeNode{I: "node2"}, &testutils.FakeNode{I: "node3"}},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node2"}, F: &testutils.FakeNode{I: "node1"}},
				&Edge{T: &testutils.FakeNode{I: "node1"}, F: &testutils.FakeNode{I: "node3"}},
			},
			input:  "node2",
			expIDs: []string{"node1"},
		},
		{
			name:  "Success/MultipleNodesFound",
			nodes: []model.Node{&testutils.FakeNode{I: "node1"}, &testutils.FakeNode{I: "node2"}, &testutils.FakeNode{I: "node3"}},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node1"}, F: &testutils.FakeNode{I: "node3"}},
				&Edge{T: &testutils.FakeNode{I: "node1"}, F: &testutils.FakeNode{I: "node2"}},
			},
			input:  "node1",
			expIDs: []string{"node2", "node3"},
		},
		{
			name:  "Success/NoNodesFound",
			nodes: []model.Node{&testutils.FakeNode{I: "node1"}, &testutils.FakeNode{I: "node2"}, &testutils.FakeNode{I: "node3"}},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node2"}, F: &testutils.FakeNode{I: "node1"}},
				&Edge{T: &testutils.FakeNode{I: "node1"}, F: &testutils.FakeNode{I: "node3"}},
			},
			input:  "node3",
			expIDs: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			collection := makeTestCollection(t, c.nodes, c.edges)
			nodes := collection.To(c.input)
			var actual []string
			for _, node := range nodes {
				actual = append(actual, node.ID())
			}
			sort.Strings(actual)
			require.Equal(t, c.expIDs, actual)
		})
	}
}

func TestCollection_From(t *testing.T) {
	type spec struct {
		name   string
		nodes  []model.Node
		edges  []model.Edge
		input  string
		expIDs []string
	}

	cases := []spec{
		{
			name:  "Success/OneNodeFound",
			nodes: []model.Node{&testutils.FakeNode{I: "node1"}, &testutils.FakeNode{I: "node2"}, &testutils.FakeNode{I: "node3"}},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node2"}, F: &testutils.FakeNode{I: "node1"}},
				&Edge{T: &testutils.FakeNode{I: "node1"}, F: &testutils.FakeNode{I: "node2"}},
			},
			input:  "node1",
			expIDs: []string{"node2"},
		},
		{
			name:  "Success/MultipleNodesFound",
			nodes: []model.Node{&testutils.FakeNode{I: "node1"}, &testutils.FakeNode{I: "node2"}, &testutils.FakeNode{I: "node3"}},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node1"}, F: &testutils.FakeNode{I: "node3"}},
				&Edge{T: &testutils.FakeNode{I: "node2"}, F: &testutils.FakeNode{I: "node3"}},
			},
			input:  "node3",
			expIDs: []string{"node1", "node2"},
		},
		{
			name:  "Success/NoNodesFound",
			nodes: []model.Node{&testutils.FakeNode{I: "node1"}, &testutils.FakeNode{I: "node2"}, &testutils.FakeNode{I: "node3"}},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node2"}, F: &testutils.FakeNode{I: "node1"}},
				&Edge{T: &testutils.FakeNode{I: "node1"}, F: &testutils.FakeNode{I: "node3"}},
			},
			input:  "node2",
			expIDs: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			collection := makeTestCollection(t, c.nodes, c.edges)
			nodes := collection.From(c.input)
			var actual []string
			for _, node := range nodes {
				actual = append(actual, node.ID())
			}
			sort.Strings(actual)
			require.Equal(t, c.expIDs, actual)
		})
	}
}

func TestCollection_Attributes(t *testing.T) {
	type spec struct {
		name          string
		nodes         []model.Node
		edges         []model.Edge
		expAttributes string
	}

	cases := []spec{
		{
			name: "Success/WithSubNodes",
			nodes: []model.Node{
				&testutils.FakeNode{
					I: "node1",
					A: attributes.Attributes{
						"title": attributes.NewString("title", "node1"),
					},
				},
				&testutils.FakeNode{
					I: "node2",
					A: attributes.Attributes{
						"title": attributes.NewString("title", "node2"),
					},
				},
				&testutils.FakeNode{
					I: "node3",
					A: attributes.Attributes{
						"title": attributes.NewString("title", "node3"),
					},
				},
			},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node2"}, F: &testutils.FakeNode{I: "node1"}},
				&Edge{T: &testutils.FakeNode{I: "node1"}, F: &testutils.FakeNode{I: "node3"}},
			},
			expAttributes: "{\"title\":\"node3\"}",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			collection := makeTestCollection(t, c.nodes, c.edges)
			attr := collection.Attributes()
			if attr == nil {
				require.Len(t, c.expAttributes, 0)
			} else {
				attrJSON, err := attr.MarshalJSON()
				require.NoError(t, err)
				require.Equal(t, c.expAttributes, string(attrJSON))
			}
		})
	}
}

func makeTestCollection(t *testing.T, nodes []model.Node, edges []model.Edge) Collection {
	c := New("test")
	for _, node := range nodes {
		require.NoError(t, c.AddNode(node))
	}
	for _, edge := range edges {
		require.NoError(t, c.AddEdge(edge))
	}
	return *c
}
