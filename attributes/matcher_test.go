package attributes

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uor-framework/client/util/testutils"
)

func TestPartialMatcher_String(t *testing.T) {
	expString := `kind=jpg,name=fish.jpg`
	m := PartialAttributeMatcher{
		"kind": "jpg",
		"name": "fish.jpg",
	}
	require.Equal(t, expString, m.String())
}

func TestPartialMatches(t *testing.T) {
	mockAttributes := testutils.MockAttributes{
		"kind":    "jpg",
		"name":    "fish.jpg",
		"another": "attribute",
	}

	n := &testutils.MockNode{A: mockAttributes}
	m := PartialAttributeMatcher{"kind": "jpg"}
	require.True(t, m.Matches(n))
}
