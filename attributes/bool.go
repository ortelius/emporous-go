package attributes

import (
	"github.com/emporous/emporous-go/model"
)

type boolAttribute bool

// NewBool returns a boolean attribute.
func NewBool(value bool) model.AttributeValue {
	return boolAttribute(value)
}

// Kind returns the kind for the attribute.
func (a boolAttribute) Kind() model.Kind {
	return model.KindBool
}

// IsNull returns whether the value is null.
func (a boolAttribute) IsNull() bool {
	return false
}

// AsBool returns the value as a boolean and errors if that is not
// the underlying type.
func (a boolAttribute) AsBool() (bool, error) {
	return bool(a), nil
}

// AsString returns the value as a string and errors if that is not
// the underlying type.
func (a boolAttribute) AsString() (string, error) {
	return "", ErrWrongKind
}

// AsFloat returns the value as a float64 value and errors if that is not
// the underlying type.
func (a boolAttribute) AsFloat() (float64, error) {
	return 0, ErrWrongKind
}

// AsInt returns the value as an int64 value and errors if that is not
// the underlying type.
func (a boolAttribute) AsInt() (int64, error) {
	return 0, ErrWrongKind
}

// AsList returns the value as a slice and errors if that is not the
// underlying type.
func (a boolAttribute) AsList() ([]model.AttributeValue, error) {
	return nil, ErrWrongKind
}

// AsObject returns the value as a map and errors if that is not the
// underlying type.
func (a boolAttribute) AsObject() (map[string]model.AttributeValue, error) {
	return nil, ErrWrongKind
}

// AsAny returns the value as an interface.
func (a boolAttribute) AsAny() interface{} {
	return bool(a)
}
