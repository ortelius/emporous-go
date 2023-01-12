package attributes

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/emporous/emporous-go/model"
)

// ErrWrongKind defines a type error try to cast an attributes value
// as the wrong type.
var ErrWrongKind = errors.New("wrong value kind")

// Attributes implements the model.Attributes interface.
type Attributes map[string]model.Attribute

var _ model.AttributeSet = &Attributes{}

// Find returns all values stored for a specified key.
func (a Attributes) Find(key string) model.Attribute {
	val, exists := a[key]
	if !exists {
		return nil
	}
	return val
}

// Exists returns whether a key,value pair exists in the
// attribute set.
func (a Attributes) Exists(input model.Attribute) (bool, error) {
	// Fail fast. Just check that the key exists and the Kinds match.
	val, ok := a[input.Key()]
	if !ok {
		return false, nil
	}

	if val.Kind() != input.Kind() {
		return false, nil
	}

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
		if val.IsNull() {
			return true, nil
		}
		return false, nil
	default:
		return false, nil
	}
}

// MarshalJSON returns a JSON formatted string representation of the
// attribute set.
func (a Attributes) MarshalJSON() ([]byte, error) {
	j := map[string]interface{}{}
	for key, value := range a {
		j[key] = value.AsAny()
	}
	jsonBytes, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

// List will list all key, value pairs for the attributes in a
// consumable format.
func (a Attributes) List() map[string]model.Attribute {
	return a
}

// Len returns the length of the attribute set.
func (a Attributes) Len() int {
	return len(a)
}

// Merge attempts to merge multiple attribute sets. If a duplicate key
// is found while merging, an error will be thrown if the value kind is
// not the same. If the value types are the same, the first set will take
// precedent.
func Merge(sets ...model.AttributeSet) (model.AttributeSet, error) {
	newSet := Attributes{}

	if len(sets) == 0 {
		return newSet, nil
	}

	if len(sets) == 1 {
		return sets[0], nil
	}

	for _, set := range sets {
		for key, value := range set.List() {
			existingVal, exists := newSet[key]
			if exists && existingVal.Kind() != value.Kind() {
				return newSet, fmt.Errorf("key %s: %w", key, ErrWrongKind)
			}
			newSet[key] = value
		}
	}

	return newSet, nil
}
