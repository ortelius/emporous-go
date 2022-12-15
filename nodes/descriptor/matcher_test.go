package descriptor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/emporous/emporous-go/attributes"
	"github.com/emporous/emporous-go/model"
	"github.com/emporous/emporous-go/util/testutils"

)

func TestJSONSubsetMatcher_Matches(t *testing.T) {
	mockAttributes := attributes.NewSet(map[string]model.AttributeValue{
		"kind":    attributes.NewString("jpg"),
		"name":    attributes.NewString("fish.jpg"),
		"another": attributes.NewString("attribute"),
	})

	n := &testutils.FakeNode{A: mockAttributes}
	m := JSONSubsetMatcher(`{"name":"fish.jpg"}`)
	match, err := m.Matches(n)
	require.NoError(t, err)
	require.True(t, match)
}
