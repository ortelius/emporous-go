package v1alpha1

import (
	"errors"
	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
)

// AttributeQueryKind object kind of AttributeQuery
const AttributeQueryKind = "AttributeQuery"

// AttributeQuery configures an attribute query.
type AttributeQuery struct {
	Kind       string `mapstructure:"kind"`
	APIVersion string `mapstructure:"apiVersion"`
	// Attributes list the configuration for Attribute types.
	Attributes []Attribute `mapstructure:"attributes"`
}

// Attribute construct a query for an individual attribute.
type Attribute struct {
	// Key represent the attribute key.
	Key string `mapstructure:"key"`
	// Value represent an attribute value.
	Value interface{} `mapstructure:"value"`
}

// ToModel converts a attribute query to a model.Attribute type.
func (a *Attribute) ToModel() (model.Attribute, error) {
	switch val := a.Value.(type) {
	case string:
		return attributes.NewString(a.Key, val), nil
	case float64:
		return attributes.NewNumber(a.Key, val), nil
	case nil:
		return attributes.NewNull(a.Key), nil
	case bool:
		return attributes.NewBool(a.Key, val), nil
	default:
		return nil, errors.New("invalid attribute type")
	}
}
