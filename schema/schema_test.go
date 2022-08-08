package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromTypes(t *testing.T) {
	type spec struct {
		name      string
		types     map[string]Type
		expSchema string
		expError  string
	}

	cases := []spec{
		{
			name: "Success/ValidConfiguration",
			types: map[string]Type{
				"test": TypeString,
				"size": TypeNumber,
			},
			expSchema: "{\"size\":{\"type\":\"number\"},\"test\":{\"type\":\"string\"}}",
		},
		{
			name: "Failure/InvalidType",
			types: map[string]Type{
				"test": TypeString,
				"size": TypeInvalid,
			},
			expError: "must set schema type",
		},
		{
			name: "Failure/UnknownType",
			types: map[string]Type{
				"test": TypeString,
				"size": 20,
			},
			expError: "unknown schema type",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			schema, err := FromTypes(c.types)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expSchema, string(schema.Export()))
			}
		})
	}
}

func TestFromBytes(t *testing.T) {
	type spec struct {
		name      string
		input     string
		expSchema string
		expError  string
	}

	cases := []spec{
		{
			name:      "Success/ValidConfiguration",
			input:     `{"size":{"type":"number"}}`,
			expSchema: "{\"size\":{\"type\":\"number\"}}",
		},
		{
			name:     "Failure/InvalidJSON",
			input:    `"size"": "type"": number`,
			expError: "schema is invalid",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			schema, err := FromBytes([]byte(c.input))
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expSchema, string(schema.Export()))
			}
		})
	}
}

func TestExport(t *testing.T) {
	exp := `{"type":"string"}`
	m := map[string]interface{}{"type": "string"}
	j, err := json.Marshal(m)
	require.NoError(t, err)
	s, err := FromBytes(j)
	require.NoError(t, err)
	require.Equal(t, exp, string(s.Export()))
}
