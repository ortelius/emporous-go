package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO(jpower432): Add more test cases for multi-level workspaces
func TestConvertFilenameForGoTemplateValue(t *testing.T) {
	type spec struct {
		name  string
		input string
		exp   string
	}

	cases := []spec{
		{
			name:  "Valid/RelativePathSameDir",
			input: "fish.jpg",
			exp:   "fish_jpg",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := ConvertFilenameForGoTemplateValue(c.input)
			require.Equal(t, c.exp, actual)
		})
	}
}
