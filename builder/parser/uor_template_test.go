package parser

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLinkableData_UOR(t *testing.T) {
	type spec struct {
		name        string
		input       string
		p           *uorParser
		tFuncs      []TemplatingFunc
		expLinksLen int
		expLinks    map[string]interface{}
		expError    string
	}

	cases := []spec{
		{
			name:        "Success/NoTemplateFunc",
			input:       "testdata/test.txt.uor",
			p:           &uorParser{filename: "test.txt.uor"},
			expLinksLen: 1,
			expLinks: map[string]interface{}{
				"fish_jpg": "fish.jpg",
			},
		},
		{
			name:  "Failure/InvalidCharacterForTemplate",
			input: "testdata/invalid.txt.uor",
			tFuncs: []TemplatingFunc{
				func(i interface{}) bool { return true },
			},
			p:        &uorParser{filename: "invalid.txt.uor"},
			expError: "template: invalid.txt.uor:1: bad character U+0026 '&'",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data, err := ioutil.ReadFile(c.input)
			require.NoError(t, err)
			if c.tFuncs != nil {
				c.p.AddFuncs(c.tFuncs...)
			}
			_, links, err := c.p.GetLinkableData(data)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {

				require.NoError(t, err)
				require.Len(t, links, c.expLinksLen)
				require.Equal(t, c.expLinks, links)
			}

		})
	}
}
