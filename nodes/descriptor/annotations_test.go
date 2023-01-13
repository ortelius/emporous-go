package descriptor

import (
	"encoding/json"
	"testing"

	empspec "github.com/emporous/collection-spec/specs-go/v1alpha1"
	"github.com/stretchr/testify/require"

	"github.com/emporous/emporous-go/attributes"
)

func TestAnnotationsFromAttributeSet(t *testing.T) {
	expMap := map[string]string{
		empspec.AnnotationEmporousAttributes: "{\"name\":\"test\",\"size\":2}",
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
		"ref":                                "example",
		empspec.AnnotationEmporousAttributes: `{"kind":"jpg","name":"fish.jpg","size":2}`,
	}
	set, err := AnnotationsToAttributeSet(annotations, nil)
	require.NoError(t, err)
	setJSON, err := set.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, expJSON, string(setJSON))
	// JSON standard lib will unmarshal all numbers as float64
	exists, err := set.Exists(attributes.NewFloat("size", 2))
	require.NoError(t, err)
	require.True(t, exists)
}

func TestAnnotationsToAttributes(t *testing.T) {
	annotations := map[string]string{
		empspec.AnnotationEmporousAttributes: "{\"name\":\"test\",\"size\":2}",
	}
	expAttrs := map[string]json.RawMessage{
		"name": []byte("\"test\""),
		"size": []byte("2"),
	}
	attrs, err := AnnotationsToAttributes(annotations)
	require.NoError(t, err)
	require.Equal(t, expAttrs, attrs)
}

func TestAnnotationsFromAttributes(t *testing.T) {
	expMap := map[string]string{
		empspec.AnnotationEmporousAttributes: "{\"name\":\"test\",\"size\":2}",
	}
	attrs := map[string]json.RawMessage{
		"name": []byte("\"test\""),
		"size": []byte("2"),
	}
	annotations, err := AnnotationsFromAttributes(attrs)
	require.NoError(t, err)
	require.Equal(t, expMap, annotations)
}

func TestAttributesFromAttributeSet(t *testing.T) {
	expAttrs := map[string]json.RawMessage{
		"name": []byte("\"test\""),
		"size": []byte("2"),
	}
	set := attributes.Attributes{
		"name": attributes.NewString("name", "test"),
		"size": attributes.NewInt("size", 2),
	}
	attrs, err := AttributesFromAttributeSet(set)
	require.NoError(t, err)
	require.Equal(t, expAttrs, attrs)
}

func TestAttributesToAttributeSet(t *testing.T) {
	expJSON := `{"test":{"kind":"jpg","name":"fish.jpg","size":2}}`
	attrs := map[string]json.RawMessage{

		"test": []byte(`{"kind":"jpg","name":"fish.jpg","size":2}`),
	}
	set, err := AttributesToAttributeSet(attrs)
	require.NoError(t, err)
	setJSON, err := set.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, expJSON, string(setJSON))
	// JSON standard lib will unmarshal all numbers as float64
	exists, err := set.Exists(attributes.NewFloat("size", 2))
	require.NoError(t, err)
	require.True(t, exists)
}
