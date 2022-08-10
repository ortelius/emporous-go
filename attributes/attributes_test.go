package attributes

import (
	"github.com/stretchr/testify/require"
	"github.com/uor-framework/uor-client-go/model"
	"testing"
)

func TestAttributes_AsJSON(t *testing.T) {
	expString := `{"name":"test","size":2}`
	test := Attributes{
		"name": NewString("name", "test"),
		"size": NewInt("size", 2),
	}
	require.Equal(t, expString, string(test.AsJSON()))
}

func TestAttributes_Exists(t *testing.T) {
	test := Attributes{
		"name": NewString("name", "test"),
		"size": NewInt("size", 2),
	}
	exists, err := test.Exists(NewString("name", "test"))
	require.NoError(t, err)
	require.True(t, exists)
	exists, err = test.Exists(NewInt("size", 2))
	require.NoError(t, err)
	require.True(t, exists)
}

func TestAttributes_Find(t *testing.T) {
	test := Attributes{
		"name": NewString("name", "test"),
		"size": NewInt("size", 2),
	}
	val := test.Find("name")
	require.Equal(t, "name", val.Key())
	require.Equal(t, model.KindString, val.Kind())
	s, err := val.AsString()
	require.NoError(t, err)
	require.Equal(t, "test", s)
}

func TestAttributes_Len(t *testing.T) {
	test := Attributes{
		"name": NewString("name", "test"),
		"size": NewInt("size", 2),
	}
	require.Equal(t, 2, test.Len())
}

func TestAttributes_List(t *testing.T) {
	test := Attributes{
		"name": NewString("name", "test"),
		"size": NewInt("size", 2),
	}
	list := test.List()
	require.Len(t, list, 2)
}
