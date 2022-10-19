package attributes

import "github.com/uor-framework/uor-client-go/model"

type sliceAttribute struct {
	key   string
	value []model.AttributeValue
}

var _ model.Attribute = sliceAttribute{}

// NewSlice returns a new string attribute.
func NewSlice(key string, attributes []model.AttributeValue) model.Attribute {
	return sliceAttribute{key: key, value: attributes}
}

// Kind returns the kind for the attribute.
func (a sliceAttribute) Kind() model.Kind {
	return model.KindString
}

// Key return the attribute key.
func (a sliceAttribute) Key() string {
	return a.key
}

// IsNull returns whether the value is null.
func (a sliceAttribute) IsNull() bool {
	return false
}

// AsBool returns the value as a boolean and errors if that is not
// the underlying type.
func (a sliceAttribute) AsBool() (bool, error) {
	return false, ErrWrongKind
}

// AsString returns the value as a string and errors if that is not
// the underlying type.
func (a sliceAttribute) AsString() (string, error) {
	return "", ErrWrongKind
}

// AsFloat returns the value as a float64 value and errors if that is not
// the underlying type.
func (a sliceAttribute) AsFloat() (float64, error) {
	return 0, ErrWrongKind
}

// AsInt returns the value as an int64 value and errors if that is not
// the underlying type.
func (a sliceAttribute) AsInt() (int64, error) {
	return 0, ErrWrongKind
}

// AsList returns the value as a slice and errors if that is not the
// underlying type.
func (a sliceAttribute) AsList() ([]model.AttributeValue, error) {
	return a.value, nil
}

// AsObject returns the value as a map and errors if that is not the
// underlying type.
func (a sliceAttribute) AsObject() (map[string]model.AttributeValue, error) {
	return nil, ErrWrongKind
}

// AsAny returns the value as an interface.
func (a sliceAttribute) AsAny() interface{} {
	return a.value
}
