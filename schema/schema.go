package schema

import (
	"errors"
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"

	"github.com/uor-framework/uor-client-go/model"
)

const (
	UnknownSchemaID   = "unknown"
	ConvertedSchemaID = "converted"
)

// Schema representation of properties in a JSON Schema format.
type Schema struct {
	jsonSchema *gojsonschema.Schema
}

// Validate performs schema validation against the
// input attribute set.
func (s *Schema) Validate(set model.AttributeSet) (bool, error) {
	attrDoc, err := set.MarshalJSON()
	if err != nil {
		return false, err
	}
	doc := gojsonschema.NewBytesLoader(attrDoc)
	result, err := s.jsonSchema.Validate(doc)
	if err != nil {
		return false, err
	}
	return result.Valid(), aggregateErrors(result.Errors())
}

func aggregateErrors(errs []gojsonschema.ResultError) error {
	if len(errs) == 0 {
		return nil
	}
	finalErr := errors.New(strings.ToLower(errs[0].String()))
	for i := 1; i < len(errs); i++ {
		finalErr = fmt.Errorf("%v:%v", finalErr, strings.ToLower(errs[i].String()))
	}
	return finalErr
}

// New create a schema from a Loader
func New(schemaLoader Loader) (Schema, error) {
	schema, err := gojsonschema.NewSchema(schemaLoader.loader)
	if err != nil {
		return Schema{}, err
	}
	return Schema{
		jsonSchema: schema,
	}, nil
}

// NewWithMulti creates a multi-schema with a root and additional loaders.
func NewWithMulti(rootSchema Loader, additionalSchemas ...Loader) (Schema, error) {
	sl := gojsonschema.NewSchemaLoader()
	for _, schema := range additionalSchemas {
		if err := sl.AddSchemas(schema.loader); err != nil {
			return Schema{}, err
		}
	}
	schema, err := sl.Compile(rootSchema.loader)
	if err != nil {
		return Schema{}, err
	}
	return Schema{
		jsonSchema: schema,
	}, nil
}
