package parser

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLinkableData_JSON(t *testing.T) {
	type spec struct {
		name        string
		input       string
		p           *jsonParser
		tFuncs      []TemplatingFunc
		expLinksLen int
		expLinks    map[string]interface{}
		expError    string
	}

	cases := []spec{
		{
			name:        "Success/NoTemplateFunc",
			input:       "testdata/test.json",
			p:           &jsonParser{filename: "test.json"},
			expLinksLen: 0,
			expLinks:    map[string]interface{}{},
		},
		{
			name:  "Success/ReplaceAllValues",
			input: "testdata/test.json",
			tFuncs: []TemplatingFunc{
				func(i interface{}) bool { return true },
			},
			p:           &jsonParser{filename: "test.json"},
			expLinksLen: 1,
			expLinks:    map[string]interface{}{"fish_jpg": "fish.jpg"},
		},
		{
			name:  "Success/ReplaceValueOnCondition",
			input: "testdata/testtag.json",
			tFuncs: []TemplatingFunc{
				func(i interface{}) bool {
					value := i.(string)
					return strings.Contains(value, "uor")
				},
			},
			p:           &jsonParser{filename: "testtag.json"},
			expLinksLen: 1,
			expLinks:    map[string]interface{}{"info_json_uor": "info.json.uor"},
		},
		{
			name:  "Failure/InvalidCharacterForTemplate",
			input: "testdata/invalid.json",
			tFuncs: []TemplatingFunc{
				func(i interface{}) bool { return true },
			},
			p:        &jsonParser{filename: "invalid.json"},
			expError: "template: invalid.json:2: bad character U+0026 '&'",
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
