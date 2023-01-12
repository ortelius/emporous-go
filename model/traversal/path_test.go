package traversal

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/emporous/emporous-go/model"
	"github.com/emporous/emporous-go/util/testutils"
)

func TestAdd(t *testing.T) {
	prev := &testutils.FakeNode{I: "node1"}
	curr := &testutils.FakeNode{I: "node2"}
	currNew := &testutils.FakeNode{I: "node3"}
	root := &testutils.FakeNode{I: "root"}
	p := NewPath(root)
	p.Add(root, prev)
	p.Add(prev, curr)

	require.Equal(t, root.ID(), p.prev[prev.I])
	require.Equal(t, prev.ID(), p.prev[curr.I])
	require.Equal(t, "", p.prev[currNew.I])
}

func TestPrev(t *testing.T) {
	prev := &testutils.FakeNode{I: "node1"}
	curr := &testutils.FakeNode{I: "node2"}
	currNew := &testutils.FakeNode{I: "node3"}
	root := &testutils.FakeNode{I: "root"}
	p := NewPath(root)
	p.Add(root, prev)
	p.Add(prev, curr)

	require.Equal(t, root, p.Prev(prev))
	require.Equal(t, prev, p.Prev(curr))
	require.Equal(t, nil, p.Prev(currNew))
}

func TestLen(t *testing.T) {
	prev := &testutils.FakeNode{I: "node1"}
	curr := &testutils.FakeNode{I: "node2"}
	currNew := &testutils.FakeNode{I: "node3"}
	root := &testutils.FakeNode{I: "root"}
	p := NewPath(root)
	p.Add(root, prev)
	p.Add(prev, curr)
	require.Equal(t, 3, p.Len())
	p.Add(curr, currNew)
	require.Equal(t, 4, p.Len())
}

func TestList(t *testing.T) {
	prev := &testutils.FakeNode{I: "node1"}
	curr := &testutils.FakeNode{I: "node2"}
	currNew := &testutils.FakeNode{I: "node3"}
	root := &testutils.FakeNode{I: "root"}
	p := NewPath(root)
	p.Add(root, prev)
	p.Add(prev, curr)
	p.Add(curr, currNew)

	require.Equal(t, []model.Node{root, prev, curr, currNew}, p.List(currNew))
}
