package matchers

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/util/testutils"
)

func TestPartialMatches(t *testing.T) {
	mockAttributes := attributes.Attributes{
		"kind":    attributes.NewString("kind", "jpg"),
		"name":    attributes.NewString("name", "fish.jpg"),
		"another": attributes.NewString("another", "attribute"),
	}

	n := &testutils.FakeNode{A: mockAttributes}
	m := PartialAttributeMatcher{"name": attributes.NewString("name", "fish.jpg")}
	match, err := m.Matches(n)
	require.NoError(t, err)
	require.True(t, match)
}
