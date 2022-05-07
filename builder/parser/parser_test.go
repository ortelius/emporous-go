package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
		{
			name:  "Valid/RelativePathNewDirecotry",
			input: "images/fish.jpg",
			exp:   "images_fish_jpg",
		},
		{
			name:  "Valid/RelativePathPreviousDirectory",
			input: "../fish.jpg",
			exp:   "___fish_jpg",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := ConvertFilenameForGoTemplateValue(c.input)
			require.Equal(t, c.exp, actual)
		})
	}
}

func TestByExtension(t *testing.T) {
	type spec struct {
		name     string
		input    string
		exp      Parser
		expError string
	}

	cases := []spec{
		{
			name:  "Success/JSON",
			input: "test.json",
			exp:   &jsonParser{filename: "test.json"},
		},
		{
			name:     "Failure/InvalidFormat",
			input:    "fish.jpg",
			expError: "format unsupported for filename: fish.jpg",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := ByExtension(c.input)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.Equal(t, c.exp, actual)
			}

		})
	}
}
