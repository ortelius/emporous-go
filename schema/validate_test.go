package schema

import (
	"github.com/stretchr/testify/require"
	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
	"testing"
)

func TestSchema_Validate(t *testing.T) {
	type spec struct {
		name     string
		schema   string
		doc      model.AttributeSet
		expRes   bool
		expError string
	}

	cases := []spec{
		{
			name:   "Success/ValidAttributes",
			schema: `{"size":{"type":"number"}}`,
			doc: attributes.Attributes{
				"size": attributes.NewFloat("size", 1.0),
			},
			expRes: true,
		},
		{
			name:   "Failure/IncompatibleType",
			schema: `{"size":{"type":"string"}}`,
			doc: attributes.Attributes{
				"size": attributes.NewFloat("size", 1.0),
			},
			expRes: false,
		},
		{
			name:   "Failure/MissingKey",
			schema: `{"size":{"type":"string"}}`,
			doc: attributes.Attributes{
				"name": attributes.NewString("name", "test"),
			},
			expRes: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			schema, err := FromBytes([]byte(c.schema))
			result, err := schema.Validate(c.doc)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expRes, result)
			}
		})
	}
}
