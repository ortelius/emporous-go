package config

import (
	"github.com/stretchr/testify/require"
	"github.com/uor-framework/uor-client-go/builder/api/v1alpha1"
	"testing"
)

func TestReadAttributeQuery(t *testing.T) {
	type spec struct {
		name     string
		path     string
		exp      v1alpha1.AttributeQuery
		expError string
	}

	cases := []spec{
		{
			name: "Success/ValidConfig",
			path: "testdata/valid-attr.yaml",
			exp: v1alpha1.AttributeQuery{
				Kind:       v1alpha1.AttributeQueryKind,
				APIVersion: v1alpha1.GroupVersion,
				Attributes: map[string]interface{}{
					"size": "small",
				},
			},
		},
		{
			name:     "Failure/InvalidConfig",
			path:     "testdata/valid-ds.yaml",
			expError: "config kind not recognized: DataSetConfiguration",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg, err := ReadAttributeQuery(c.path)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.exp, cfg)
			}
		})
	}
}

func TestReadCollectionConfig(t *testing.T) {
	type spec struct {
		name     string
		path     string
		exp      v1alpha1.DataSetConfiguration
		expError string
	}

	cases := []spec{
		{
			name: "Success/ValidConfig",
			path: "testdata/valid-ds.yaml",
			exp: v1alpha1.DataSetConfiguration{
				Kind:       v1alpha1.DataSetConfigurationKind,
				APIVersion: v1alpha1.GroupVersion,
				Collection: v1alpha1.Collection{
					SchemaAddress: "localhost:5001/schema:latest",
					Files: []v1alpha1.File{
						{
							File: "*.json",
							Attributes: map[string]interface{}{
								"fiction": true,
							},
						},
					},
				},
			},
		},
		{
			name:     "Failure/InvalidConfig",
			path:     "testdata/valid-attr.yaml",
			expError: "config kind not recognized: AttributeQuery",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg, err := ReadCollectionConfig(c.path)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.exp, cfg)
			}
		})
	}
}
