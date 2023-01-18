package config

import (
	"testing"

	empspec "github.com/emporous/collection-spec/specs-go/v1alpha1"
	"github.com/stretchr/testify/require"

	"github.com/emporous/emporous-go/api/client/v1alpha1"
	"github.com/emporous/emporous-go/schema"
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
				TypeMeta: v1alpha1.TypeMeta{
					Kind:       v1alpha1.AttributeQueryKind,
					APIVersion: v1alpha1.GroupVersion,
				},
				Attributes: []byte(`{"size":"small"}`),
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
		exp      v1alpha1.DataSetConfiguration
		expError string
	}

	cases := []spec{
		{
			name: "Success/ValidConfig",
			path: "testdata/valid-ds.yaml",
			exp: v1alpha1.DataSetConfiguration{
				TypeMeta: v1alpha1.TypeMeta{
					Kind:       v1alpha1.DataSetConfigurationKind,
					APIVersion: v1alpha1.GroupVersion,
				},
				Collection: v1alpha1.DataSetConfigurationSpec{
					SchemaAddress: "localhost:5001/schema:latest",
					Files: []v1alpha1.File{
						{
							File: "*.json",
							Attributes: map[string]interface{}{
								"fiction": true,
							},
							FileInfo: empspec.File{
								UID: -1,
								GID: -1,
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

func TestReadSchemaConfiguration(t *testing.T) {
	type spec struct {
		name     string
		path     string
		exp      v1alpha1.SchemaConfiguration
		expError string
	}

	cases := []spec{
		{
			name: "Success/ValidConfig",
			path: "testdata/valid-schema.yaml",
			exp: v1alpha1.SchemaConfiguration{
				TypeMeta: v1alpha1.TypeMeta{
					Kind:       v1alpha1.SchemaConfigurationKind,
					APIVersion: v1alpha1.GroupVersion,
				},
				Schema: v1alpha1.SchemaConfigurationSpec{
					AttributeTypes: map[string]schema.Type{
						"test": schema.TypeString,
					},
				},
			},
		},
		{
			name:     "Failure/InvalidConfig",
			path:     "testdata/valid-attr.yaml",
			expError: "config kind AttributeQuery, does not match expected SchemaConfiguration",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg, err := ReadSchemaConfig(c.path)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.exp, cfg)
			}
		})
	}
}
