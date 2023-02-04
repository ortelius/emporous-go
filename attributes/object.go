package attributes

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/emporous/emporous-go/model"
)

// ErrWrongKind defines a type error try to cast an attributes value
// as the wrong type.
var ErrWrongKind = errors.New("wrong value kind")

type mapAttribute map[string]model.AttributeValue

var _ model.AttributeValue = mapAttribute{}
var _ model.AttributeSet = mapAttribute{}

// NewObject returns a new object attribute.
func NewObject(attributes map[string]model.AttributeValue) model.AttributeValue {
	return mapAttribute(attributes)
}

// NewSet returns an object as an attribute set.
func NewSet(attributes map[string]model.AttributeValue) model.AttributeSet {
	return mapAttribute(attributes)
}

// Kind returns the kind for the attribute.
func (a mapAttribute) Kind() model.Kind {
	return model.KindObject
}

// IsNull returns whether the value is null.
func (a mapAttribute) IsNull() bool {
	return false
}

// AsBool returns the value as a boolean and errors if that is not
// the underlying type.
func (a mapAttribute) AsBool() (bool, error) {
	return false, ErrWrongKind
}

// AsString returns the value as a string and errors if that is not
// the underlying type.
func (a mapAttribute) AsString() (string, error) {
	return "", ErrWrongKind
}

// AsFloat returns the value as a float64 value and errors if that is not
// the underlying type.
func (a mapAttribute) AsFloat() (float64, error) {
	return 0, ErrWrongKind
}

// AsInt returns the value as an int64 value and errors if that is not
// the underlying type.
func (a mapAttribute) AsInt() (int64, error) {
	return 0, ErrWrongKind
}

// AsList returns the value as a slice and errors if that is not the
// underlying type.
func (a mapAttribute) AsList() ([]model.AttributeValue, error) {
	return nil, ErrWrongKind
}

// AsObject returns the value as a map and errors if that is not the
// underlying type.
func (a mapAttribute) AsObject() (map[string]model.AttributeValue, error) {
	return a, nil
}

// AsAny returns the value as an interface.
func (a mapAttribute) AsAny() interface{} {
	return a
}

// Find returns the attribute values stored for a specified key.
// This only support top-level keys and does not descend into
// child objects.
func (a mapAttribute) Find(key string) model.AttributeValue {
	val, exists := a[key]
	if !exists {
		return nil
	}
	return val
}

// Exists returns whether a key,value pair exists in the
// attribute set.
func (a mapAttribute) Exists(key string, value model.AttributeValue) (bool, error) {
	// Fail fast. Just check that the key exists and the Kinds match.
	setVal, ok := a[key]
	if !ok {
		return false, nil
	}

	if setVal.Kind() != value.Kind() {
		return false, nil
	}

	return checkValue(value, setVal)
}

func checkObject(inputVal, setVal model.AttributeValue) (bool, error) {
	inputMap, err := inputVal.AsObject()
	if err != nil {
		return false, err
	}

	valMap, err := setVal.AsObject()
	if err != nil {
		return false, err
	}

	for key, value := range inputMap {
		otherVal, ok := valMap[key]
		if !ok {
			return false, nil
		}
		match, err := checkValue(value, otherVal)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}

	return true, nil
}

func checkList(input, val model.AttributeValue) (bool, error) {
	inputList, err := input.AsList()
	if err != nil {
		return false, err
	}
	valList, err := val.AsList()
	if err != nil {
		return false, err
	}

	// FIXME(jpower432): This could be an issue because of sorting.
	// We really just want to check the the presence of the same attributes. M
	// Most json diff libraries I have found are strict about array ordering.
	// We don't need to be here.
	for i := 0; i < len(inputList); i++ {
		match, err := checkValue(inputList[i], valList[i])
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}

	return true, nil
}

func checkValue(input, val model.AttributeValue) (bool, error) {
	switch input.Kind() {
	case model.KindString:
		outS, err := val.AsString()
		if err != nil {
			return false, err
		}
		inS, err := input.AsString()
		if err != nil {
			return false, err
		}
		return outS == inS, nil
	case model.KindFloat:
		outF, err := val.AsFloat()
		if err != nil {
			return false, err
		}
		inF, err := input.AsFloat()
		if err != nil {
			return false, err
		}
		return outF == inF, nil
	case model.KindInt:
		outI, err := val.AsInt()
		if err != nil {
			return false, err
		}
		inI, err := input.AsInt()
		if err != nil {
			return false, err
		}
		return outI == inI, nil
	case model.KindBool:
		outB, err := val.AsBool()
		if err != nil {
			return false, err
		}
		inB, err := input.AsBool()
		if err != nil {
			return false, err
		}
		return outB == inB, nil
	case model.KindNull:
		return val.IsNull(), nil
	case model.KindObject:
		return checkObject(input, val)
	case model.KindList:
		return checkList(input, val)
	default:
		return false, fmt.Errorf("unsupported type")
	}
}

// MarshalJSON returns a JSON formatted string representation of the
// attribute set.
func (a mapAttribute) MarshalJSON() ([]byte, error) {
	alias := map[string]model.AttributeValue{}
	for k, v := range a {
		alias[k] = v
	}
	return json.Marshal(alias)
}

// List will list all key, value pairs for the attributes in a
// consumable format.
func (a mapAttribute) List() map[string]model.AttributeValue {
	return a
}

// Len returns the length of the attribute set.
func (a mapAttribute) Len() int {
	return len(a)
}

// MergeOptions define options for attribute merging.
type MergeOptions struct {
	AllowSameTypeOverwrites bool
}

// Merge complete a two-way merge of multiple sets of attributes. If AllowSameTypeOverwrites is set to true, patches can overwrite
// values of the same type (default false). If false, the function will return an error  when a key collision occurs.
func Merge(original map[string]model.AttributeValue, opts MergeOptions, patches ...map[string]model.AttributeValue) (map[string]model.AttributeValue, error) {
	if len(patches) == 0 {
		return original, nil
	}

	var err error
	currSet := NewObject(original)
	for _, patch := range patches {
		patchObject := NewObject(patch)
		currSet, err = mergeObjects(currSet, patchObject, nil, opts)
		if err != nil {
			return nil, err
		}
	}

	mergedList := currSet.(mapAttribute)
	return mergedList, nil
}

func mergeValue(path []string, patch model.AttributeValue, key string, value model.AttributeValue, opts MergeOptions) (model.AttributeValue, error) {
	patchObject, err := patch.AsObject()
	if err != nil {
		return nil, err
	}
	patchValue, patchHasValue := patchObject[key]

	if !patchHasValue {
		return value, nil
	}

	patchValueIsObject := patchValue.Kind() == model.KindObject

	path = append(path, key)
	pathStr := strings.Join(path, ".")

	if value.Kind() == model.KindObject {
		if !patchValueIsObject {
			return value, fmt.Errorf("patch value must be object for key %q", pathStr)
		}

		return mergeObjects(value, patchValue, path, opts)
	}

	if value.Kind() == model.KindList && patchValueIsObject {
		return mergeObjects(value, patchValue, path, opts)
	}

	if value.Kind() != patchValue.Kind() {
		return nil, fmt.Errorf("path %q: %w", pathStr, ErrWrongKind)
	}

	if !opts.AllowSameTypeOverwrites {
		match, err := checkValue(value, patchValue)
		if err != nil {
			return nil, fmt.Errorf("path %q: %w", pathStr, err)
		}
		if !match {
			return nil, fmt.Errorf("cannot overwrite value at %q", pathStr)
		}
	}

	return patchValue, nil
}

func mergeObjects(data, patch model.AttributeValue, path []string, opts MergeOptions) (model.AttributeValue, error) {
	if patch.Kind() == model.KindObject {
		if data.Kind() == model.KindList {
			dataArray, err := data.AsList()
			if err != nil {
				return nil, err
			}

			ret := make([]model.AttributeValue, len(dataArray))

			for i, val := range dataArray {
				ret[i], err = mergeValue(path, patch, strconv.Itoa(i), val, opts)
				if err != nil {
					return NewList(ret), err
				}
			}

			return NewList(ret), nil
		} else if data.Kind() == model.KindObject {
			dataObject, err := data.AsObject()
			if err != nil {
				return nil, err
			}
			ret := mapAttribute{}

			for k, v := range dataObject {
				ret[k], err = mergeValue(path, patch, k, v, opts)
				if err != nil {
					return ret, err
				}
			}

			patchObject, err := patch.AsObject()
			if err != nil {
				return nil, err
			}
			// Add in new objects from patches
			for key, value := range patchObject {
				_, ok := ret[key]
				if !ok {
					ret[key] = value
				}
			}

			return ret, nil
		}
	}

	return data, nil
}
