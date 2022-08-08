package attributes

import (
	"github.com/stretchr/testify/require"
	"github.com/uor-framework/uor-client-go/model"
	"testing"
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

func TestBoolAttribute_AsNumber(t *testing.T) {
	test := NewBool("test", false)
	b, err := test.AsNumber()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, float64(0), b)
}

func TestBoolAttribute_AsString(t *testing.T) {
	test := NewBool("test", false)
	b, err := test.AsString()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, "", b)
}
