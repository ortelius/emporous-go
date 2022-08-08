package attributes

import "github.com/uor-framework/uor-client-go/model"

type nullAttribute struct {
	key string
}

var _ model.Attribute = nullAttribute{}

// NewNull returns a null attribute.
func NewNull(key string) model.Attribute {
	return nullAttribute{key: key}
}

// Kind returns the kind for the attribute.
func (a nullAttribute) Kind() model.Kind {
	return model.KindNull
}

// Key return the attribute key.
func (a nullAttribute) Key() string {
	return a.key
}

// IsNull returns whether the value is null.
func (a nullAttribute) IsNull() bool {
	return true
}

// AsBool returns the value as a boolean errors if that is not
// the underlying type.
func (a nullAttribute) AsBool() (bool, error) {
	return false, ErrWrongKind
}

// AsString returns the value as a string errors if that is not
// the underlying type.
func (a nullAttribute) AsString() (string, error) {
	return "", ErrWrongKind
}

// AsNumber returns the value as a number value errors if that is not
// the underlying type.
func (a nullAttribute) AsNumber() (float64, error) {
	return 0, ErrWrongKind
}

// AsAny returns the value as an interface.
func (a nullAttribute) AsAny() interface{} {
	return nil
}
