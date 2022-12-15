package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/emporous/emporous-go/api/client/v1alpha1"
	"github.com/emporous/emporous-go/attributes"
	"github.com/emporous/emporous-go/model"
)

func TestConvertToModel(t *testing.T) {
	type spec struct {
		name       string
		attributes v1alpha1.Attributes
		asserFunc  func(set model.AttributeSet) bool
		expError   string
	}

	cases := []spec{
		{
			name: "Success/OneAttributeKind",
			attributes: v1alpha1.Attributes{
				"test": 2.0,
				"size": 5.2,
			},
			asserFunc: func(set model.AttributeSet) bool {
				testExists, err := set.Exists("test", attributes.NewFloat(2.0))
				if err != nil {
					return false
				}
				sizeExists, err := set.Exists("size", attributes.NewFloat(5.2))
				if err != nil {
					return false
				}
				return testExists && sizeExists
			},
		},
		{
			name: "Success/MultipleAttributeKinds",
			attributes: v1alpha1.Attributes{
				"test":     "a test",
				"istest":   true,
				"other":    nil,
				"size":     5.2,
				"sequence": 1,
			},
			asserFunc: func(set model.AttributeSet) bool {
				stringExists, err := set.Exists("test", attributes.NewString("a test"))
				if err != nil {
					t.Log(err)
					return false
				}
				boolExists, err := set.Exists("istest", attributes.NewBool(true))
				if err != nil {
					t.Log(err)
					return false
				}
				nullExists, err := set.Exists("other", attributes.NewNull())
				if err != nil {
					t.Log(err)
					return false
				}
				numExists, err := set.Exists("size", attributes.NewFloat(5.2))
				if err != nil {
					t.Log(err)
					return false
				}
				intExists, err := set.Exists("sequence", attributes.NewInt(1))
				if err != nil {
					t.Log(err)
					return false
				}
				return stringExists && boolExists && numExists && nullExists && intExists
			},
		},
		{
			name: "Failure/InvalidAttributeType",
			attributes: v1alpha1.Attributes{
				"test": struct{}{},
			},
			expError: "error converting attribute test to model: invalid attribute type",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mdl, err := ConvertToModel(c.attributes)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.True(t, c.asserFunc(mdl))
			}
		})
	}
}
