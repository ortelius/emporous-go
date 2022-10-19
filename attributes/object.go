package attributes

import "github.com/uor-framework/uor-client-go/model"

type mapAttribute struct {
	key   string
	value map[string]model.AttributeValue
}

var _ model.Attribute = mapAttribute{}

// NewMap returns a new string attribute.
func NewMap(key string, attributes map[string]model.AttributeValue) model.Attribute {
	return mapAttribute{key: key, value: attributes}
}

// Kind returns the kind for the attribute.
func (a mapAttribute) Kind() model.Kind {
	return model.KindString
}

// Key return the attribute key.
func (a mapAttribute) Key() string {
	return a.key
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
	return a.value, nil
}

// AsAny returns the value as an interface.
func (a mapAttribute) AsAny() interface{} {
	return a.value
}
