package attributes

import (
	"github.com/stretchr/testify/require"
	"github.com/uor-framework/uor-client-go/model"
	"testing"
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

func TestStringAttribute_AsNumber(t *testing.T) {
	test := NewBool("test", false)
	s, err := test.AsNumber()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, float64(0), s)
}

func TestStringAttribute_AsString(t *testing.T) {
	test := NewString("test", "testvalue")
	s, err := test.AsString()
	require.NoError(t, err)
	require.Equal(t, "testvalue", s)
}
