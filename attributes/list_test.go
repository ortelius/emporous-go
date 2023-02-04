package attributes

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/emporous/emporous-go/model"
)

func TestListAttribute_Kind(t *testing.T) {
	test := NewList([]model.AttributeValue{
		NewString("testvalue"),
	})
	require.Equal(t, model.KindList, test.Kind())
}

func TestListAttribute_AsBool(t *testing.T) {
	test := NewList([]model.AttributeValue{
		NewString("testvalue"),
	})
	s, err := test.AsBool()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, false, s)
}

func TestListAttribute_AsFloat(t *testing.T) {
	test := NewList([]model.AttributeValue{
		NewString("testvalue"),
	})
	s, err := test.AsFloat()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, float64(0), s)
}

func TestListAttribute_AsInt(t *testing.T) {
	test := NewList([]model.AttributeValue{
		NewString("testvalue"),
	})
	s, err := test.AsInt()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, int64(0), s)
}

func TestListAttribute_AsString(t *testing.T) {
	test := NewList([]model.AttributeValue{})
	s, err := test.AsString()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, "", s)
}

func TestListAttribute_IsNull(t *testing.T) {
	test := NewList([]model.AttributeValue{
		NewString("testvalue"),
	})
	require.False(t, test.IsNull())
}

func TestListAttribute_AsList(t *testing.T) {
	test := NewList([]model.AttributeValue{
		NewString("testvalue"),
	})
	s, err := test.AsList()
	require.NoError(t, err)
	require.Equal(t, []model.AttributeValue{NewString("testvalue")}, s)

}

func TestListAttribute_AsObject(t *testing.T) {
	test := NewList([]model.AttributeValue{
		NewString("testvalue"),
	})
	s, err := test.AsObject()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, map[string]model.AttributeValue(nil), s)
}
