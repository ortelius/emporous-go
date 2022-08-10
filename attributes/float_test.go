package attributes

import (
	"github.com/stretchr/testify/require"
	"github.com/uor-framework/uor-client-go/model"
	"testing"
)

func TestFloatAttribute_Kind(t *testing.T) {
	test := NewFloat("test", 1)
	require.Equal(t, model.KindFloat, test.Kind())
}

func TestFloatAttribute_AsBool(t *testing.T) {
	test := NewFloat("test", 1)
	n, err := test.AsBool()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, false, n)
}

func TestFloatAttribute_AsFloat(t *testing.T) {
	test := NewFloat("test", 1)
	n, err := test.AsFloat()
	require.NoError(t, err)
	require.Equal(t, float64(1), n)
}

func TestFloatAttribute_AsInt(t *testing.T) {
	test := NewFloat("test", 1)
	n, err := test.AsInt()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, 0, n)
}

func TestFloatAttribute_AsString(t *testing.T) {
	test := NewFloat("test", 1)
	n, err := test.AsString()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, "", n)
}

func TestFloatAttribute_IsNull(t *testing.T) {
	test := NewFloat("test", 1.0)
	require.False(t, test.IsNull())
}
