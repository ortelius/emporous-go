package attributes

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/emporous/emporous-go/model"
)

func TestBoolAttribute_Kind(t *testing.T) {
	test := NewBool("test", true)
	require.Equal(t, model.KindBool, test.Kind())
}

func TestBoolAttribute_AsBool(t *testing.T) {
	test := NewBool("test", true)
	b, err := test.AsBool()
	require.NoError(t, err)
	require.Equal(t, true, b)
}

func TestBoolAttribute_AsFloat(t *testing.T) {
	test := NewBool("test", false)
	b, err := test.AsFloat()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, float64(0), b)
}

func TestBoolAttribute_AsInt(t *testing.T) {
	test := NewBool("test", false)
	b, err := test.AsInt()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, int64(0), b)
}

func TestBoolAttribute_AsString(t *testing.T) {
	test := NewBool("test", false)
	b, err := test.AsString()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, "", b)
}

func TestBoolAttribute_IsNull(t *testing.T) {
	test := NewBool("test", false)
	require.False(t, test.IsNull())
}
