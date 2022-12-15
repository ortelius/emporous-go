package collection

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/emporous/emporous-go/model"
	"github.com/emporous/emporous-go/util/testutils"
	"github.com/emporous/emporous-go/attributes"
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
				A: attributes.NewSet(map[string]model.AttributeValue{
					"kind": attributes.NewString("txt"),
					"name": attributes.NewString("test"),
				}),
			},
		},
		want: []model.Node{
			&testutils.FakeNode{
				I: "node1",
				A: attributes.NewSet(map[string]model.AttributeValue{
					"kind": attributes.NewString("txt"),
					"name": attributes.NewString("test"),
				}),
			},
		},
	},
	{
		nodes: []model.Node{
			&testutils.FakeNode{
				I: "node1",
				A: attributes.NewSet(map[string]model.AttributeValue{
					"kind": attributes.NewString("txt"),
					"name": attributes.NewString("test"),
				}),
			},
			&testutils.FakeNode{
				I: "node2",
				A: attributes.NewSet(map[string]model.AttributeValue{
					"kind": attributes.NewString("txt"),
				}),
			},
		},
		want: []model.Node{
			&testutils.FakeNode{
				I: "node2",
				A: attributes.NewSet(map[string]model.AttributeValue{
					"kind": attributes.NewString("txt"),
				}),
			},
			&testutils.FakeNode{
				I: "node1",
				A: attributes.NewSet(map[string]model.AttributeValue{
					"kind": attributes.NewString("txt"),
					"name": attributes.NewString("test"),
				}),
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
