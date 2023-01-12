package schema

import (
	"encoding/json"
	"sort"

	"github.com/xeipuuv/gojsonschema"
)

type Loader struct {
	loader gojsonschema.JSONLoader
	raw    json.RawMessage
}

// Export returns the json raw message.
func (l Loader) Export() json.RawMessage {
	return l.raw
}

// FromTypes builds a JSON Schema from a key with an associated type.
// All keys provided will be considered required types in the schema when
// comparing sets of attributes.
func FromTypes(types Types) (Loader, error) {
	if err := types.Validate(); err != nil {
		return Loader{}, err
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
		return Loader{}, err
	}
	return FromBytes(b)
}

// FromGo loads a Go struct into a JSON schema that
// can be used for attribute validation.
func FromGo(source interface{}) (Loader, error) {
	loader := gojsonschema.NewGoLoader(source)

	raw, err := loader.LoadJSON()
	if err != nil {
		return Loader{}, err
	}
	return Loader{
		loader: loader,
		raw:    raw.([]byte),
	}, nil
}

// FromBytes loads data into a JSON Schema that can be used
// for attribute validation.
func FromBytes(data []byte) (Loader, error) {
	loader := gojsonschema.NewBytesLoader(data)
	return Loader{
		loader: loader,
		raw:    data,
	}, nil
}
