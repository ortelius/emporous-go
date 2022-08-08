package matchers

import (
	"github.com/uor-framework/uor-client-go/attributes"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uor-framework/uor-client-go/util/testutils"
)

func TestPartialMatches(t *testing.T) {
	mockAttributes := attributes.Attributes{
		"kind":    attributes.NewString("kind", "jpg"),
		"name":    attributes.NewString("name", "fish.jpg"),
		"another": attributes.NewString("another", "attribute"),
	}

	n := &testutils.MockNode{A: mockAttributes}
	m := PartialAttributeMatcher{"name": attributes.NewString("name", "fish.jpg")}
	require.True(t, m.Matches(n))
}
