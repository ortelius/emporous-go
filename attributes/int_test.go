package attributes

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/emporous/emporous-go/model"
)

func TestIntAttribute_Kind(t *testing.T) {
	test := NewInt(1)
	require.Equal(t, model.KindInt, test.Kind())
}

func TestIntAttribute_AsBool(t *testing.T) {
	test := NewInt(1)
	n, err := test.AsBool()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, false, n)
}

func TestIntAttribute_AsInt(t *testing.T) {
	test := NewInt(1)
	n, err := test.AsInt()
	require.NoError(t, err)
	require.Equal(t, int64(1), n)
}

func TestIntAttribute_AsFloat(t *testing.T) {
	test := NewInt(1)
	n, err := test.AsFloat()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, float64(0), n)
}

func TestIntAttribute_AsString(t *testing.T) {
	test := NewInt(1)
	n, err := test.AsString()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, "", n)
}

func TestIntAttribute_IsNull(t *testing.T) {
	test := NewInt(1)
	require.False(t, test.IsNull())
}
