package schema

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/emporous/emporous-go/model"
)

// Type represent types for an attribute schema.
type Type int

const (
	TypeInvalid Type = iota
	TypeNull
	TypeBool
	TypeNumber
	TypeInteger
	TypeString
)

// String prints a string representation of the attribute kind.
func (t Type) String() string {
	return stringByType[t]
}

// IsLike returns the model Kind that correlate to the schema Type.
func (t Type) IsLike() (model.Kind, error) {
	// Use d to represent the default kind if one does not match
	var d model.Kind
	if err := t.validate(); err != nil {
		return d, err
	}
	kind, found := modelKindByType[t]
	if !found {
		return d, fmt.Errorf("type %s is not linked to a model kind", t.String())
	}
	return kind, nil
}

// UnmarshalJSON unmarshal a JSON serialized type to the Schema Type
func (t *Type) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	*t = typeByString[j]
	return t.validate()
}

// MarshalJSON marshals the Schema Type into JSON format.
func (t Type) MarshalJSON() ([]byte, error) {
	if err := t.validate(); err != nil {
		return nil, err
	}
	return json.Marshal(t.String())
}

// validate performs basic validation
// on a Type.
func (t Type) validate() error {
	if _, found := stringByType[t]; found {
		return nil
	}
	switch t {
	case TypeInvalid:
		// TypeInvalid is the default value for the concrete type, which means the field was not set.
		return errors.New("must set schema type")
	default:
		return fmt.Errorf("unknown schema type")
	}
}

// stringByType maps the schema Type to its string
// representation.
var stringByType = map[Type]string{
	TypeNumber:  "number",
	TypeInteger: "integer",
	TypeBool:    "boolean",
	TypeString:  "string",
	TypeNull:    "null",
}

// typeByString maps the string representation of the schema Type
// to the schema Type.
var typeByString = map[string]Type{
	"number":  TypeNumber,
	"integer": TypeInteger,
	"boolean": TypeBool,
	"string":  TypeString,
	"null":    TypeNull,
}

// modelKindByType maps each schema type to a
// corresponding model Kind.
var modelKindByType = map[Type]model.Kind{
	TypeNumber:  model.KindFloat,
	TypeInteger: model.KindInt,
	TypeBool:    model.KindBool,
	TypeString:  model.KindString,
	TypeNull:    model.KindNull,
	TypeInvalid: model.KindInvalid,
}

// Types represent a schema Type mapped to a key of string type.
type Types map[string]Type

// Validate performs basic validation
// on a set of schema types.
func (t Types) Validate() error {
	for _, value := range t {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}
