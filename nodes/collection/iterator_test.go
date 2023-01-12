package collection

import (
	"testing"

	"github.com/uor-framework/uor-client-go/attributes"

	"github.com/stretchr/testify/require"

	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/util/testutils"
)

var iteratorTests = []struct {
	nodes []model.Node
	want  []model.Node
}{
	{nodes: nil, want: nil},
	{
		nodes: []model.Node{
			&testutils.FakeNode{
				I: "node1",
				A: attributes.Attributes{
					"kind": attributes.NewString("kind", "txt"),
					"name": attributes.NewString("name", "test"),
				},
			},
		},
		want: []model.Node{
			&testutils.FakeNode{
				I: "node1",
				A: attributes.Attributes{
					"kind": attributes.NewString("kind", "txt"),
					"name": attributes.NewString("name", "test"),
				},
			},
		},
	},
	{
		nodes: []model.Node{
			&testutils.FakeNode{
				I: "node1",
				A: attributes.Attributes{
					"kind": attributes.NewString("kind", "txt"),
					"name": attributes.NewString("name", "test"),
				},
			},
			&testutils.FakeNode{
				I: "node2",
				A: attributes.Attributes{
					"kind": attributes.NewString("kind", "txt"),
				},
			},
		},
		want: []model.Node{
			&testutils.FakeNode{
				I: "node2",
				A: attributes.Attributes{
					"kind": attributes.NewString("kind", "txt"),
				},
			},
			&testutils.FakeNode{
				I: "node1",
				A: attributes.Attributes{
					"kind": attributes.NewString("kind", "txt"),
					"name": attributes.NewString("name", "test"),
				},
			},
		},
	},
}

func TestByAttributeIterator(t *testing.T) {
	for _, test := range iteratorTests {
		it := NewByAttributesIterator(test.nodes)
		for i := 0; i < 2; i++ {
			require.Equal(t, it.Len(), len(test.nodes))
			var got []model.Node
			for it.Next() {
				got = append(got, it.Node())
				require.Equal(t, len(got)+it.Len(), len(test.nodes))
			}
			require.Equal(t, test.want, got)
			it.Reset()
		}
	}
}
