package schema

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/xeipuuv/gojsonschema"
)

// Schema representation of properties in a JSON Schema format.
type Schema struct {
	JSONSchema *gojsonschema.Schema
	raw        json.RawMessage
}

// Export returns the json raw message.
func (s Schema) Export() json.RawMessage {
	return s.raw
}

// FromTypes builds a JSON Schema from a key with an associated type.
// All keys provided will be considered required types in the schema when
// comparing sets of attributes.
func FromTypes(types Types) (Schema, error) {
	if err := types.Validate(); err != nil {
		return Schema{}, err
	}

	// Build an object in json from the provided types
	type jsonSchema struct {
		Type       string                       `json:"type"`
		Properties map[string]map[string]string `json:"properties"`
		Required   []string                     `json:"required"`
	}

	// Fill in properties and required keys. At this point
	// we consider all keys as required.
	properties := map[string]map[string]string{}
	var required []string
	for key, value := range types {
		properties[key] = map[string]string{"type": value.String()}
		required = append(required, key)
	}

	// Make the required slice order deterministic
	sort.Slice(required, func(i, j int) bool {
		return required[i] < required[j]
	})

	tmp := jsonSchema{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}
	b, err := json.Marshal(tmp)
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
		return Schema{}, fmt.Errorf("error creating JSON schema: %w", err)
	}
	return Schema{
		JSONSchema: schema,
		raw:        data,
	}, nil
}
