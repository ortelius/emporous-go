package attributes

import (
	"errors"

	"github.com/uor-framework/uor-client-go/model"
)

// Reflect will create a model.Attribute type from a Go type.
func Reflect(key string, value interface{}) (model.Attribute, error) {
	switch typVal := value.(type) {
	case string:
		return NewString(key, typVal), nil
	case float64:
		return NewNumber(key, typVal), nil
	case nil:
		return NewNull(key), nil
	case bool:
		return NewBool(key, typVal), nil
	default:
		return nil, errors.New("invalid attribute type")
	}
}
