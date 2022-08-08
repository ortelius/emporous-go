package schema

import (
	"encoding/json"
	"errors"
	"github.com/uor-framework/uor-client-go/model"
)

// Type represent the Attribute Kinds
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
	switch t {
	case TypeInvalid:
		return "INVALID"
	case TypeNull:
		return "null"
	case TypeBool:
		return "bool"
	case TypeInteger:
		return "integer"
	case TypeNumber:
		return "number"
	case TypeString:
		return "string"
	default:
		return ""
	}
}

// IsLike returns the model Kind that correlate to the schema Type.
func (t Type) IsLike() (model.Kind, error) {
	switch t {
	case TypeInvalid:
		return model.KindInvalid, nil
	case TypeNull:
		return model.KindNull, nil
	case TypeBool:
		return model.KindBool, nil
	case TypeInteger:
		return model.KindInt, nil
	case TypeNumber:
		return model.KindFloat, nil
	case TypeString:
		return model.KindString, nil
	default:
		return model.KindInvalid, errors.New(" schema type")
	}
}

// UnmarshalJSON unmarshal a JSON serialized type to the Schema Type
func (t *Type) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	typ, err := getType(j)
	if err != nil {
		return err
	}
	*t = typ
	return nil
}

// MarshalJSON marshals the Schema Type into JSON format.
func (t Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// validate performs basic validation
// on a Type.
func (t Type) validate() error {
	typ, err := getType(t.String())
	if err != nil {
		return err
	}
	if typ == TypeInvalid {
		return errors.New("must set schema type")
	}
	return nil
}

// getType will return a schema type based on a string.
// FIXME(jpower432): Do this with maps?
func getType(s string) (Type, error) {
	switch s {
	case "bool":
		return TypeBool, nil
	case "null":
		return TypeNull, nil
	case "number":
		return TypeNumber, nil
	case "string":
		return TypeString, nil
	case "INVALID":
		return TypeInvalid, nil
	default:
		return TypeInvalid, errors.New("unknown schema type")
	}
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
