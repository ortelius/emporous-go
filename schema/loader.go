package schema

import (
	"encoding/json"

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
