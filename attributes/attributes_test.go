package attributes

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/emporous/emporous-go/model"
)

func TestAttributes_MarshalJSON(t *testing.T) {
	expString := `{"name":"test","size":2}`
	test := Attributes{
		"name": NewString("name", "test"),
		"size": NewInt("size", 2),
	}
	testJSON, err := test.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, expString, string(testJSON))
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

func TestMerge(t *testing.T) {
	type spec struct {
		name      string
		set1      Attributes
		set2      Attributes
		expString string
		expError  string
	}

	cases := []spec{
		{
			name: "Success/MergedAttributes",
			set1: Attributes{
				"name": NewString("name", "snoopy"),
				"size": NewInt("size", 2),
			},
			set2: Attributes{
				"breed": NewString("breed", "beagle"),
			},
			expString: `{"breed":"beagle","name":"snoopy","size":2}`,
		},
		{
			name: "Success/MergedAttributesOverwrite",
			set1: Attributes{
				"name": NewString("name", "snoopy"),
				"size": NewInt("size", 2),
			},
			set2: Attributes{
				"name":  NewString("name", "pluto"),
				"breed": NewString("breed", "beagle"),
			},
			expString: `{"breed":"beagle","name":"pluto","size":2}`,
		},
		{
			name: "Failure/TypeMismatch",
			set1: Attributes{
				"name": NewString("name", "snoopy"),
				"size": NewInt("size", 2),
			},
			set2: Attributes{
				"breed": NewString("breed", "beagle"),
				"size":  NewString("size", "medium"),
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
				testJSON, err := mergedSet.MarshalJSON()
				require.NoError(t, err)
				require.Equal(t, c.expString, string(testJSON))
			}

		})
	}
}
