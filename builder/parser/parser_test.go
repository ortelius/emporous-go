package parser

import (
	"io/ioutil"
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
			name:  "Valid/RelativePathNewDirectory",
			input: "images/fish.jpg",
			exp:   "images_fish_jpg",
		},
		{
			name:  "Valid/SpecialCharacter",
			input: "images/fish-1.jpg",
			exp:   "images_fish_1_jpg",
		},
		{
			name:  "Valid/SpecialCharacterAccepted",
			input: "images/fish-1.jpg",
			exp:   "images_fish_1_jpg",
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

func TestContentType(t *testing.T) {
	type spec struct {
		name     string
		input    string
		exp      Parser
		expError string
	}

	cases := []spec{
		{
			name:  "Success/JSON",
			input: "testdata/test.json",
			exp:   &jsonParser{filename: "testdata/test.json"},
		},
		{
			name:     "Failure/InvalidFormat",
			input:    "testdata/fish.jpg",
			expError: "format unsupported for filename: testdata/fish.jpg",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data, err := ioutil.ReadFile(c.input)
			require.NoError(t, err)
			actual, err := ByContentType(c.input, data)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.Equal(t, c.exp, actual)
			}

		})
	}
}
