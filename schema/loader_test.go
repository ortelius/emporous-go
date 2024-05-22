package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

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
