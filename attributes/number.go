package attributes

import "github.com/uor-framework/uor-client-go/model"

type numberAttribute struct {
	key   string
	value float64
}

var _ model.Attribute = numberAttribute{}

// NewNumber returns a number attribute.
func NewNumber(key string, value float64) model.Attribute {
	return numberAttribute{key: key, value: value}
}

// Kind returns the kind for the attribute.
func (a numberAttribute) Kind() model.Kind {
	return model.KindNumber
}

// Key return the attribute key.
func (a numberAttribute) Key() string {
	return a.key
}

// IsNull returns whether the value is null.
func (a numberAttribute) IsNull() bool {
	return false
}

// AsBool returns the value as a boolean errors if that is not
// the underlying type.
func (a numberAttribute) AsBool() (bool, error) {
	return false, ErrWrongKind
}

// AsString returns the value as a string errors if that is not
// the underlying type.
func (a numberAttribute) AsString() (string, error) {
	return "", ErrWrongKind
}

// AsNumber returns the value as a number value errors if that is not
// the underlying type.
func (a numberAttribute) AsNumber() (float64, error) {
	return a.value, nil
}

// AsAny returns the value as an interface.
func (a numberAttribute) AsAny() interface{} {
	return a.value
}
