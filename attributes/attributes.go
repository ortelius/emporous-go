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
func (a Attributes) Exists(key string, kind model.Kind, value interface{}) bool {
	val, ok := a[key]
	if !ok {
		return false
	}

	if val.Kind() != kind {
		return false
	}

	switch kind {
	case model.KindString:
		s, err := val.AsString()
		if err != nil {
			return false
		}
		return s == value.(string)
	case model.KindNumber:
		n, err := val.AsNumber()
		if err != nil {
			return false
		}
		return n == value.(float64)
	case model.KindBool:
		b, err := val.AsBool()
		if err != nil {
			return false
		}
		return b == value.(bool)
	case model.KindNull:
		if val.IsNull() {
			return true
		}
		fallthrough
	default:
		return false
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
