package attributes

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/emporous/emporous-go/model"
)

func TestObjectAttribute_Kind(t *testing.T) {
	test := NewObject(map[string]model.AttributeValue{
		"key": NewString("testvalue"),
	})
	require.Equal(t, model.KindObject, test.Kind())
}

func TestObjectAttribute_AsBool(t *testing.T) {
	test := NewObject(map[string]model.AttributeValue{
		"key": NewString("testvalue"),
	})
	s, err := test.AsBool()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, false, s)
}

func TestObjectAttribute_AsFloat(t *testing.T) {
	test := NewObject(map[string]model.AttributeValue{
		"key": NewString("testvalue"),
	})
	s, err := test.AsFloat()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, float64(0), s)
}

func TestObjectAttribute_AsInt(t *testing.T) {
	test := NewObject(map[string]model.AttributeValue{
		"key": NewString("testvalue"),
	})
	s, err := test.AsInt()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, int64(0), s)
}

func TestObjectAttribute_AsString(t *testing.T) {
	test := NewObject(map[string]model.AttributeValue{
		"key": NewString("testvalue"),
	})
	s, err := test.AsString()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, "", s)
}

func TestObjectAttribute_IsNull(t *testing.T) {
	test := NewObject(map[string]model.AttributeValue{
		"key": NewString("testvalue"),
	})
	require.False(t, test.IsNull())
}

func TestObjectAttribute_AsList(t *testing.T) {
	test := NewObject(map[string]model.AttributeValue{
		"key": NewString("testvalue"),
	})
	s, err := test.AsList()
	require.ErrorIs(t, ErrWrongKind, err)
	require.Equal(t, []model.AttributeValue(nil), s)
}

func TestObjectAttribute_AsObject(t *testing.T) {
	test := NewObject(map[string]model.AttributeValue{
		"key": NewString("testvalue"),
	})
	s, err := test.AsObject()
	require.NoError(t, err)
	require.Equal(t, map[string]model.AttributeValue{"key": NewString("testvalue")}, s)
}

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
		set       mapAttribute
		patches   []map[string]model.AttributeValue
		mergeOpts MergeOptions
		expString string
		expError  string
	}

	cases := []spec{
		{
			name: "Success/MergedAttributes",
			set: mapAttribute{
				"name": NewString("snoopy"),
				"size": NewInt(2),
			},
			patches: []map[string]model.AttributeValue{
				{"breed": NewString("beagle")},
			},
			expString: `{"breed":"beagle","name":"snoopy","size":2}`,
		},
		{
			name: "Success/MergedAttributesOverwrite",
			set: mapAttribute{
				"name": NewString("snoopy"),
				"size": NewInt(2),
			},
			mergeOpts: MergeOptions{
				AllowSameTypeOverwrites: true,
			},
			patches: []map[string]model.AttributeValue{
				{"name": NewString("pluto"),
					"breed": NewString("beagle")},
			},
			expString: `{"breed":"beagle","name":"pluto","size":2}`,
		},
		{
			name: "Success/MergedMultiplePatches",
			set: mapAttribute{
				"name": NewString("snoopy"),
				"size": NewInt(2),
			},
			mergeOpts: MergeOptions{
				AllowSameTypeOverwrites: true,
			},
			patches: []map[string]model.AttributeValue{
				{"name": NewString("pluto"), "breed": NewString("beagle")},
				{"name": NewString("bingo"), "breed": NewString("sheepdog"), "color": NewString("brown")},
			},
			expString: `{"breed":"sheepdog","color":"brown","name":"bingo","size":2}`,
		},
		{
			name: "Success/ObjectMerges",
			set: mapAttribute{
				"name": NewString("bingo"),
				"description": NewObject(map[string]model.AttributeValue{
					"color":      NewString("brown"),
					"brightness": NewString("dark"),
					"age":        NewInt(4),
					"owner":      NewString("farmer"),
				}),
			},
			mergeOpts: MergeOptions{
				AllowSameTypeOverwrites: true,
			},
			patches: []map[string]model.AttributeValue{
				{
					"name": NewString("bingo"),
					"description": NewObject(map[string]model.AttributeValue{
						"spelling": NewList([]model.AttributeValue{
							NewString("b"),
							NewString("i"),
							NewString("n"),
							NewString("g"),
							NewString("o"),
						}),
					}),
				},
			},
			expString: `{"description":{"age":4,"brightness":"dark","color":"brown","owner":"farmer","spelling":["b","i","n","g","o"]},"name":"bingo"}`,
		},
		{
			name: "Success/ListMerges",
			set: mapAttribute{
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
			},
			mergeOpts: MergeOptions{
				AllowSameTypeOverwrites: true,
			},
			patches: []map[string]model.AttributeValue{
				{
					"name": NewString("bingo"),
					"description": NewObject(map[string]model.AttributeValue{
						"spelling": NewObject(map[string]model.AttributeValue{
							"0": NewString(""),
						}),
					}),
				},
			},
			expString: `{"description":{"age":4,"brightness":"dark","color":"brown","owner":"farmer","spelling":["","i","n","g","o"]},"name":"bingo"}`,
		},
		{
			name: "Failure/DisallowOverwrite",
			set: mapAttribute{
				"name": NewString("snoopy"),
				"size": NewString("small"),
			},
			patches: []map[string]model.AttributeValue{
				{"breed": NewString("beagle"),
					"size": NewString("medium")},
			},
			expError: `cannot overwrite value at "size"`,
		},
		{
			name: "Failure/TypeMismatch",
			set: mapAttribute{
				"name": NewString("snoopy"),
				"size": NewInt(2),
			},
			patches: []map[string]model.AttributeValue{
				{"breed": NewString("beagle"),
					"size": NewString("medium")},
			},
			expError: `path "size": wrong value kind`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mergedSet, err := Merge(c.set, c.mergeOpts, c.patches...)
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
