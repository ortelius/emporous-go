package schema

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
)

func TestSchema_Validate(t *testing.T) {
	type spec struct {
		name        string
		schemaTypes Types
		doc         model.AttributeSet
		expRes      bool
		expError    string
	}

	cases := []spec{
		{
			name: "Success/ValidAttributes",
			schemaTypes: map[string]Type{
				"size": TypeNumber,
			},
			doc: attributes.Attributes{
				"size": attributes.NewFloat("size", 1.0),
			},
			expRes: true,
		},
		{
			name: "Failure/IncompatibleType",
			schemaTypes: map[string]Type{
				"size": TypeBool,
			},
			doc: attributes.Attributes{
				"size": attributes.NewFloat("size", 1.0),
			},
			expRes:   false,
			expError: "size: invalid type. expected: boolean, given: integer",
		},
		{
			name: "Failure/MissingKey",
			schemaTypes: map[string]Type{
				"size": TypeString,
			},
			doc: attributes.Attributes{
				"name": attributes.NewString("name", "test"),
			},
			expError: "(root): size is required",
			expRes:   false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			schema, err := FromTypes(c.schemaTypes)
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
