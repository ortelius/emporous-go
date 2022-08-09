package config

import (
	"github.com/stretchr/testify/require"
	"github.com/uor-framework/uor-client-go/builder/api/v1alpha1"
	"github.com/uor-framework/uor-client-go/model"
	"testing"
)

func TestConvertToModel(t *testing.T) {
	type spec struct {
		name      string
		query     v1alpha1.AttributeQuery
		asserFunc func(set model.AttributeSet) bool
		expError  string
	}

	cases := []spec{
		{
			name: "Success/OneAttributeKind",
			query: v1alpha1.AttributeQuery{
				Kind:       v1alpha1.AttributeQueryKind,
				APIVersion: "test",
				Attributes: []v1alpha1.Attribute{
					{
						Key:   "test",
						Value: 2.0,
					},
					{
						Key:   "size",
						Value: 5.2,
					},
				},
			},
			asserFunc: func(set model.AttributeSet) bool {
				return set.Exists("test", model.KindNumber, 2.0) && set.Exists("size", model.KindNumber, 5.2)
			},
		},
		{
			name: "Success/MultipleAttributeKinds",
			query: v1alpha1.AttributeQuery{
				Kind:       v1alpha1.AttributeQueryKind,
				APIVersion: "test",
				Attributes: []v1alpha1.Attribute{
					{
						Key:   "test",
						Value: "a test",
					},
					{
						Key:   "istest",
						Value: true,
					},
					{
						Key:   "other",
						Value: nil,
					},
					{
						Key:   "size",
						Value: 5.2,
					},
				},
			},
			asserFunc: func(set model.AttributeSet) bool {
				return (set.Exists("test", model.KindString, "a test") && set.Exists("istest", model.KindBool, true)) &&
					set.Exists("other", model.KindNull, nil) && set.Exists("size", model.KindNumber, 5.2)

			},
		},
		{
			name: "Failure/InvalidAttributeType",
			query: v1alpha1.AttributeQuery{
				Kind:       v1alpha1.AttributeQueryKind,
				APIVersion: "test",
				Attributes: []v1alpha1.Attribute{
					{
						Key:   "test",
						Value: struct{}{},
					},
				},
			},
			expError: "error converting attribute test to model: invalid attribute type",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			model, err := ConvertToModel(c.query)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.True(t, c.asserFunc(model))
			}
		})
	}
}
