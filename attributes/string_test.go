package attributes

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/uor-framework/uor-client-go/model"
)

func TestStringAttribute_Kind(t *testing.T) {
	test := NewString("test", "testvalue")
	require.Equal(t, model.KindString, test.Kind())
}

func TestStringAttribute_AsBool(t *testing.T) {
	test := NewString("test", "testvalue")
	s, err := test.AsBool()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, false, s)
}

func TestStringAttribute_AsFloat(t *testing.T) {
	test := NewBool("test", false)
	s, err := test.AsFloat()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, float64(0), s)
}

func TestStringAttribute_AsInt(t *testing.T) {
	test := NewBool("test", false)
	s, err := test.AsInt()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, int64(0), s)
}

func TestStringAttribute_AsString(t *testing.T) {
	test := NewString("test", "testvalue")
	s, err := test.AsString()
	require.NoError(t, err)
	require.Equal(t, "testvalue", s)
}

func TestStringAttribute_IsNull(t *testing.T) {
	test := NewString("test", "testvalue")
	require.False(t, test.IsNull())
}
