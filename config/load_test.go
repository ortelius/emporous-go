package config

import (
	"github.com/stretchr/testify/require"
	v1alpha12 "github.com/uor-framework/uor-client-go/api/v1alpha1"
	"testing"
)

func TestReadAttributeQuery(t *testing.T) {
	type spec struct {
		name     string
		path     string
		exp      v1alpha12.AttributeQuery
		expError string
	}

	cases := []spec{
		{
			name: "Success/ValidConfig",
			path: "testdata/valid-attr.yaml",
			exp: v1alpha12.AttributeQuery{
				TypeMeta: v1alpha12.TypeMeta{
					Kind:       v1alpha12.AttributeQueryKind,
					APIVersion: v1alpha12.GroupVersion,
				},
				Attributes: map[string]interface{}{
					"size": "small",
				},
			},
		},
		{
			name:     "Failure/InvalidConfig",
			path:     "testdata/valid-ds.yaml",
			expError: "config kind DataSetConfiguration, does not match expected AttributeQuery",
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

func TestReadDataSetConfig(t *testing.T) {
	type spec struct {
		name     string
		path     string
		exp      v1alpha12.DataSetConfiguration
		expError string
	}

	cases := []spec{
		{
			name: "Success/ValidConfig",
			path: "testdata/valid-ds.yaml",
			exp: v1alpha12.DataSetConfiguration{
				TypeMeta: v1alpha12.TypeMeta{
					Kind:       v1alpha12.DataSetConfigurationKind,
					APIVersion: v1alpha12.GroupVersion,
				},
				Collection: v1alpha12.Collection{
					SchemaAddress: "localhost:5001/schema:latest",
					Files: []v1alpha12.File{
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
			expError: "config kind AttributeQuery, does not match expected DataSetConfiguration",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg, err := ReadDataSetConfig(c.path)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.exp, cfg)
			}
		})
	}
}
