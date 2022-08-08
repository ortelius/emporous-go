package descriptor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnnotationsToAttributes(t *testing.T) {
	expJSON := `{"kind":"jpg","name":"fish.jpg"}`
	annotations := map[string]string{
		"kind": "jpg",
		"name": "fish.jpg",
	}
	require.Equal(t, expJSON, string(AnnotationsToAttributes(annotations).AsJSON()))
}
