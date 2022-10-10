package schema

import (
	"errors"
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"

	"github.com/uor-framework/uor-client-go/model"
)

// Validate performs schema validation against the
// input attribute set.
func (s Schema) Validate(set model.AttributeSet) (bool, error) {
	doc := gojsonschema.NewBytesLoader(set.AsJSON())
	result, err := s.JSONSchema.Validate(doc)
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
