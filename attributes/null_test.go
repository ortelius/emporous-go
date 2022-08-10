package attributes

import (
	"github.com/stretchr/testify/require"
	"github.com/uor-framework/uor-client-go/model"
	"testing"
)

func TestNullAttribute_Kind(t *testing.T) {
	test := NewNull("test")
	require.Equal(t, model.KindNull, test.Kind())
}

func TestNullAttribute_AsBool(t *testing.T) {
	test := NewNull("test")
	n, err := test.AsBool()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, false, n)
}

func TestNullAttribute_AsInt(t *testing.T) {
	test := NewNull("test")
	n, err := test.AsInt()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, 0, n)
}

func TestNullAttribute_AsFloat(t *testing.T) {
	test := NewNull("test")
	n, err := test.AsFloat()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, float64(0), n)
}

func TestNullAttribute_AsString(t *testing.T) {
	test := NewNull("test")
	n, err := test.AsString()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, "", n)
}

func TestNullAttribute_IsNull(t *testing.T) {
	test := NewNull("test")
	require.True(t, test.IsNull())
}
