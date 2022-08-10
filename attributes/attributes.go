package attributes

import (
	"encoding/json"
	"errors"
	"github.com/uor-framework/uor-client-go/model"
)

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

// AsJSON returns a JSON formatted string representation of the
// attribute set. If the values are not valid, nil is returned.
func (a Attributes) AsJSON() json.RawMessage {
	j := map[string]interface{}{}
	for key, value := range a {
		j[key] = value.AsAny()
	}
	jsonBytes, err := json.Marshal(j)
	if err != nil {
		return nil
	}
	return jsonBytes
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
