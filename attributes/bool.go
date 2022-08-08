package attributes

import (
	"github.com/uor-framework/uor-client-go/model"
)

type boolAttribute struct {
	key   string
	value bool
}

var _ model.Attribute = boolAttribute{}

// NewBool returns a boolean attribute.
func NewBool(key string, value bool) model.Attribute {
	return boolAttribute{key: key, value: value}
}

// Kind returns the kind for the attribute.
func (a boolAttribute) Kind() model.Kind {
	return model.KindBool
}

// Key return the attribute key.
func (a boolAttribute) Key() string {
	return a.key
}

// IsNull returns whether the value is null.
func (a boolAttribute) IsNull() bool {
	return false
}

// AsBool returns the value as a boolean errors if that is not
// the underlying type.
func (a boolAttribute) AsBool() (bool, error) {
	return a.value, nil
}

// AsString returns the value as a string errors if that is not
// the underlying type.
func (a boolAttribute) AsString() (string, error) {
	return "", ErrWrongKind
}

// AsNumber returns the value as a number value errors if that is not
// the underlying type.
func (a boolAttribute) AsNumber() (float64, error) {
	return 0, ErrWrongKind
}

// AsAny returns the value as an interface.
func (a boolAttribute) AsAny() interface{} {
	return a.value
}
