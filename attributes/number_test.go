package attributes

import (
	"github.com/stretchr/testify/require"
	"github.com/uor-framework/uor-client-go/model"
	"testing"
)

func TestNumberAttribute_Kind(t *testing.T) {
	test := NewNumber("test", 1)
	require.Equal(t, model.KindNumber, test.Kind())
}

func TestNumberAttribute_AsBool(t *testing.T) {
	test := NewNumber("test", 1)
	n, err := test.AsBool()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, false, n)
}

func TestNumberAttribute_AsNumber(t *testing.T) {
	test := NewNumber("test", 1)
	n, err := test.AsNumber()
	require.NoError(t, err)
	require.Equal(t, float64(1), n)
}

func TestNumberAttribute_AsString(t *testing.T) {
	test := NewNumber("test", 1)
	n, err := test.AsString()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, "", n)
}
