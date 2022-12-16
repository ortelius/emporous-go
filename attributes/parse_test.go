package attributes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/emporous/emporous-go/model"
)

func TestParse(t *testing.T) {
	type spec struct {
		name       string
		data       json.RawMessage
		assertFunc func(value model.AttributeValue) bool
		expError   string
	}

	cases := []spec{
		{
			name: "Success/PrimitiveValues",
			data: json.RawMessage(`"test": "test"`),
			assertFunc: func(value model.AttributeValue) bool {
				return value.Kind() == model.KindString
			},
		},
		{
			name: "Success/ObjectType",
			data: json.RawMessage(`{"test": {"test": true}}`),
			assertFunc: func(value model.AttributeValue) bool {
				return value.Kind() == model.KindObject
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			attr, err := Parse(c.data)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.True(t, c.assertFunc(attr))
			}
		})
	}
}

func TestParseToSet(t *testing.T) {
	type spec struct {
		name       string
		attributes json.RawMessage
		asserFunc  func(set model.AttributeSet) bool
		expError   string
	}

	cases := []spec{
		{
			name:       "Success/OneAttributeKind",
			attributes: json.RawMessage(`{"size":5.2,"test": 2.0}`),
			asserFunc: func(set model.AttributeSet) bool {
				testExists, err := set.Exists("test", NewFloat(2.0))
				if err != nil {
					return false
				}
				sizeExists, err := set.Exists("size", NewFloat(5.2))
				if err != nil {
					return false
				}
				return testExists && sizeExists
			},
		},
		{
			name:       "Success/MultipleAttributeKinds",
			attributes: json.RawMessage(`{"istest":true,"other":null,"size":5.2,"test":"a test"}`),
			asserFunc: func(set model.AttributeSet) bool {
				stringExists, err := set.Exists("test", NewString("a test"))
				if err != nil {
					return false
				}
				boolExists, err := set.Exists("istest", NewBool(true))
				if err != nil {
					return false
				}
				nullExists, err := set.Exists("other", NewNull())
				if err != nil {
					return false
				}
				numExists, err := set.Exists("size", NewFloat(5.2))
				if err != nil {
					return false
				}
				return stringExists && boolExists && nullExists && numExists
			},
		},
		{
			name:       "Failure/InvalidAttributeType",
			attributes: json.RawMessage(`{"size": struct{}}`),
			expError:   "invalid character 's' looking for beginning of value",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mdl, err := ParseToSet(c.attributes)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.True(t, c.asserFunc(mdl))
			}
		})
	}
}
