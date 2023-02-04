package attributes

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/emporous/emporous-go/model"
)

func TestFloatAttribute_Kind(t *testing.T) {
	test := NewFloat(1)
	require.Equal(t, model.KindFloat, test.Kind())
}

func TestFloatAttribute_AsBool(t *testing.T) {
	test := NewFloat(1)
	n, err := test.AsBool()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, false, n)
}

func TestFloatAttribute_AsFloat(t *testing.T) {
	test := NewFloat(1)
	n, err := test.AsFloat()
	require.NoError(t, err)
	require.Equal(t, float64(1), n)
}

func TestFloatAttribute_AsInt(t *testing.T) {
	test := NewFloat(1)
	n, err := test.AsInt()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, int64(0), n)
}

func TestFloatAttribute_AsString(t *testing.T) {
	test := NewFloat(1)
	n, err := test.AsString()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, "", n)
}

func TestFloatAttribute_IsNull(t *testing.T) {
	test := NewFloat(1.0)
	require.False(t, test.IsNull())
}

func TestFloatAttribute_AsList(t *testing.T) {
	test := NewFloat(1)
	s, err := test.AsList()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, []model.AttributeValue(nil), s)
}

func TestFloatAttribute_AsObject(t *testing.T) {
	test := NewFloat(1)
	s, err := test.AsObject()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, map[string]model.AttributeValue(nil), s)
}
