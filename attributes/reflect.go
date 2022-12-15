package attributes

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/emporous/emporous-go/model"
)

// ErrInvalidAttribute defines the error thrown when an attribute has an invalid
// type.
var ErrInvalidAttribute = errors.New("invalid attribute type")

// Reflect will create a model.AttributeValue type from a Go type.
func Reflect(value interface{}) (model.AttributeValue, error) {
	// Try type switch first
	switch typVal := value.(type) {
	case string:
		return NewString(typVal), nil
	case float64:
		return NewFloat(typVal), nil
	case int64:
		return NewInt(typVal), nil
	case nil:
		return NewNull(), nil
	case bool:
		return NewBool(typVal), nil
	}

	// To catch more types try reflection
	reflectVal := reflect.ValueOf(value)
	switch reflectVal.Kind() {
	case reflect.Bool:
		return NewBool(reflectVal.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewInt(reflectVal.Int()), nil
	case reflect.Float32, reflect.Float64:
		return NewFloat(reflectVal.Float()), nil
	case reflect.String:
		return NewString(reflectVal.String()), nil
	default:
		return nil, ErrInvalidAttribute
	}
}

// Parse will create a model.AttributeValue type from json.RawMessage.
func Parse(value json.RawMessage) (model.AttributeValue, error) {
	// TODO(jpower432): Finish this and use as a more performant alternative
	// to reflect when using json.
	return nil, nil
}
