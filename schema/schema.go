package schema

import (
	"encoding/json"
	"github.com/xeipuuv/gojsonschema"
)

// Schema representation of properties in a JSON Schema format.
type Schema struct {
	*gojsonschema.Schema
	raw json.RawMessage
}

// Export returns the json raw message.
func (s Schema) Export() json.RawMessage {
	return s.raw
}

// FromTypes builds a JSON Schema from a key with an associated type.
func FromTypes(types Types) (Schema, error) {
	if err := types.Validate(); err != nil {
		return Schema{}, err
	}

	properties := map[string]map[string]string{}
	for key, value := range types {
		properties[key] = map[string]string{"type": value.String()}
	}
	b, err := json.Marshal(properties)
	if err != nil {
		return Schema{}, err
	}
	return FromBytes(b)
}

// FromBytes loads data into a JSON Schema that can be used
// for attribute validation.
func FromBytes(data []byte) (Schema, error) {
	loader := gojsonschema.NewBytesLoader(data)
	schema, err := gojsonschema.NewSchema(loader)
	if err != nil {
		return Schema{}, err
	}
	return Schema{
		Schema: schema,
		raw:    data,
	}, nil
}
