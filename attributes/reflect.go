package attributes

import (
	"errors"
	"reflect"

	"github.com/uor-framework/uor-client-go/model"
)

// ErrInvalidAttribute defines the error thrown when an attribute has an invalid
// type.
var ErrInvalidAttribute = errors.New("invalid attribute type")

// Reflect will create a model.Attribute type from a Go type.
func Reflect(key string, value interface{}) (model.Attribute, error) {
	// Try type switch first
	switch typVal := value.(type) {
	case string:
		return NewString(key, typVal), nil
	case float64:
		return NewFloat(key, typVal), nil
	case int64:
		return NewInt(key, typVal), nil
	case nil:
		return NewNull(key), nil
	case bool:
		return NewBool(key, typVal), nil
	}

	// To catch more types try reflection
	reflectVal := reflect.ValueOf(value)
	switch reflectVal.Kind() {
	case reflect.Bool:
		return NewBool(key, reflectVal.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewInt(key, reflectVal.Int()), nil
	case reflect.Float32, reflect.Float64:
		return NewFloat(key, reflectVal.Float()), nil
	case reflect.String:
		return NewString(key, reflectVal.String()), nil
	default:
		return nil, ErrInvalidAttribute
	}
}
