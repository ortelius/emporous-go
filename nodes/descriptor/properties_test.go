package descriptor

import (
	"testing"

	"github.com/stretchr/testify/require"
	empspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	"github.com/emporous/emporous-go/attributes"
	"github.com/emporous/emporous-go/model"
)

func TestProperties_MarshalJSON(t *testing.T) {
	expJSON := `{"core-descriptor":{"id":"id","name":"","version":"","type":"","foundBy":"","locations":null,"licenses":null,"language":"","cpes":null,"purl":""},"core-link":{"registryHint":"test","namespaceHint":"namespace","transitive":false},"test":{"name":"test","size":2}}`
	set := attributes.Attributes{
		"name": attributes.NewString("name", "test"),
		"size": attributes.NewInt("size", 2),
	}
	props := &Properties{
		Link: &empspec.LinkAttributes{
			RegistryHint:  "test",
			NamespaceHint: "namespace",
		},
		Descriptor: &empspec.DescriptorAttributes{
			Component: empspec.Component{
				ID: "id",
			},
		},
		Others: map[string]model.AttributeSet{"test": set},
	}
	propsJSON, err := props.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, expJSON, string(propsJSON))
}
