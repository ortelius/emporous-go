package schema

import (
	"github.com/uor-framework/uor-client-go/model"
	"github.com/xeipuuv/gojsonschema"
)

// Validate performs schema validation against the
// input attribute set.
func (s Schema) Validate(set model.AttributeSet) (bool, error) {
	doc := gojsonschema.NewBytesLoader(set.AsJSON())
	result, err := s.Schema.Validate(doc)
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}
