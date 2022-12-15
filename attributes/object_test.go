package attributes

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/uor-framework/uor-client-go/model"
)

func TestAttributes_MarshalJSON(t *testing.T) {
	expString := `{"name":"test","size":2}`
	test := mapAttribute{
		"name": NewString("test"),
		"size": NewInt(2),
	}
	testJSON, err := test.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, expString, string(testJSON))
}

func TestAttributes_Exists(t *testing.T) {
	test := mapAttribute{
		"name": NewString("bingo"),
		"description": NewObject(map[string]model.AttributeValue{
			"color":      NewString("brown"),
			"brightness": NewString("dark"),
			"age":        NewInt(4),
			"owner":      NewString("farmer"),
			"spelling": NewList([]model.AttributeValue{
				NewString("b"),
				NewString("i"),
				NewString("n"),
				NewString("g"),
				NewString("o"),
			}),
		}),
	}
	exists, err := test.Exists("name", NewString("bingo"))
	require.NoError(t, err)
	require.True(t, exists)
	exists, err = test.Exists("description", NewObject(map[string]model.AttributeValue{
		"color": NewString("brown"),
	}))
	require.NoError(t, err)
	require.True(t, exists)
	exists, err = test.Exists("description", NewObject(map[string]model.AttributeValue{
		"spelling": NewList([]model.AttributeValue{
			NewString("b"),
			NewString("i"),
		}),
	}))
	require.NoError(t, err)
	require.True(t, exists)
}

func TestAttributes_Find(t *testing.T) {
	test := mapAttribute{
		"name": NewString("test"),
		"size": NewInt(2),
	}
	val := test.Find("name")
	require.Equal(t, model.KindString, val.Kind())
	s, err := val.AsString()
	require.NoError(t, err)
	require.Equal(t, "test", s)
}

func TestAttributes_Len(t *testing.T) {
	test := mapAttribute{
		"name": NewString("test"),
		"size": NewInt(2),
	}
	require.Equal(t, 2, test.Len())
}

func TestAttributes_List(t *testing.T) {
	test := mapAttribute{
		"name": NewString("test"),
		"size": NewInt(2),
	}
	list := test.List()
	require.Len(t, list, 2)
}

func TestMerge(t *testing.T) {
	type spec struct {
		name      string
		set1      mapAttribute
		set2      mapAttribute
		expString string
		expError  string
	}

	cases := []spec{
		{
			name: "Success/MergedAttributes",
			set1: mapAttribute{
				"name": NewString("snoopy"),
				"size": NewInt(2),
			},
			set2: mapAttribute{
				"breed": NewString("beagle"),
			},
			expString: `{"breed":"beagle","name":"snoopy","size":2}`,
		},
		{
			name: "Success/MergedAttributesOverwrite",
			set1: mapAttribute{
				"name": NewString("snoopy"),
				"size": NewInt(2),
			},
			set2: mapAttribute{
				"name":  NewString("pluto"),
				"breed": NewString("beagle"),
			},
			expString: `{"breed":"beagle","name":"pluto","size":2}`,
		},
		{
			name: "Failure/TypeMismatch",
			set1: mapAttribute{
				"name": NewString("snoopy"),
				"size": NewInt(2),
			},
			set2: mapAttribute{
				"breed": NewString("beagle"),
				"size":  NewString("medium"),
			},
			expError: "key size: wrong value kind",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mergedSet, err := Merge(c.set1, c.set2)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				mergedObject := NewSet(mergedSet)
				testJSON, err := mergedObject.MarshalJSON()
				require.NoError(t, err)
				require.Equal(t, c.expString, string(testJSON))
			}

		})
	}
}
