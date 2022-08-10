package ocimanifest

import (
	"github.com/stretchr/testify/require"
	"github.com/uor-framework/uor-client-go/attributes"
	"testing"
)

func TestAnnotationsFromAttributeSet(t *testing.T) {
	expMap := map[string]string{
		AnnotationUORAttributes: "{\"name\":\"test\",\"size\":2}",
	}
	set := attributes.Attributes{
		"name": attributes.NewString("name", "test"),
		"size": attributes.NewInt("size", 2),
	}
	annotations, err := AnnotationsFromAttributeSet(set)
	require.NoError(t, err)
	require.Equal(t, expMap, annotations)
}

func TestAnnotationsToAttributeSet(t *testing.T) {
	expJSON := `{"kind":"jpg","name":"fish.jpg","ref":"example","size":2}`
	annotations := map[string]string{
		"ref":                   "example",
		AnnotationUORAttributes: `{"kind":"jpg","name":"fish.jpg","size":2}`,
	}
	set, err := AnnotationsToAttributeSet(annotations, nil)
	require.NoError(t, err)
	require.Equal(t, expJSON, string(set.AsJSON()))
	// JSON standard lib will unmarshal all numbers as float64
	exists, err := set.Exists(attributes.NewFloat("size", 2))
	require.True(t, exists)
}
