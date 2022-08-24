package attributes

import "github.com/uor-framework/uor-client-go/model"

type intAttribute struct {
	key   string
	value int64
}

var _ model.Attribute = intAttribute{}

// NewInt returns an int attribute.
func NewInt(key string, value int64) model.Attribute {
	return intAttribute{key: key, value: value}
}

// Kind returns the kind for the attribute.
func (a intAttribute) Kind() model.Kind {
	return model.KindInt
}

// Key return the attribute key.
func (a intAttribute) Key() string {
	return a.key
}

// IsNull returns whether the value is null.
func (a intAttribute) IsNull() bool {
	return false
}

// AsBool returns the value as a boolean and errors if that is not
// the underlying type.
func (a intAttribute) AsBool() (bool, error) {
	return false, ErrWrongKind
}

// AsString returns the value as a string and errors if that is not
// the underlying type.
func (a intAttribute) AsString() (string, error) {
	return "", ErrWrongKind
}

// AsFloat returns the value as a float value and errors if that is not
// the underlying type.
func (a intAttribute) AsFloat() (float64, error) {
	return 0, ErrWrongKind
}

// AsInt returns the value as an int value errors and if that is not
// the underlying type.
func (a intAttribute) AsInt() (int64, error) {
	return a.value, nil
}

// AsAny returns the value as an interface.
func (a intAttribute) AsAny() interface{} {
	return a.value
}
