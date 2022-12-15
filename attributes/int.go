package attributes

import "github.com/emporous/emporous-go/model"

type intAttribute int64

// NewInt returns an int attribute.
func NewInt(value int64) model.AttributeValue {
	return intAttribute(value)
}

// Kind returns the kind for the attribute.
func (a intAttribute) Kind() model.Kind {
	return model.KindInt
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

// AsFloat returns the value as a float64 value and errors if that is not
// the underlying type.
func (a intAttribute) AsFloat() (float64, error) {
	return 0, ErrWrongKind
}

// AsInt returns the value as an int64 value and errors if that is not
// the underlying type.
func (a intAttribute) AsInt() (int64, error) {
	return int64(a), nil
}

// AsList returns the value as a slice and errors if that is not the
// underlying type.
func (a intAttribute) AsList() ([]model.AttributeValue, error) {
	return nil, ErrWrongKind
}

// AsObject returns the value as a map and errors if that is not the
// underlying type.
func (a intAttribute) AsObject() (map[string]model.AttributeValue, error) {
	return nil, ErrWrongKind
}

// AsAny returns the value as an interface.
func (a intAttribute) AsAny() interface{} {
	return int64(a)
}
