package config

import (
	"fmt"

	"github.com/emporous/emporous-go/api/client/v1alpha1"
	"github.com/emporous/emporous-go/attributes"
	"github.com/emporous/emporous-go/model"
)

// ConvertToModel converts v1alpha1.Attributes to an model.AttributeSet.
func ConvertToModel(input v1alpha1.Attributes) (model.AttributeSet, error) {
	set := map[string]model.AttributeValue{}
	for key, val := range input {
		mattr, err := attributes.Reflect(val)
		if err != nil {
			return nil, fmt.Errorf("error converting attribute %s to model: %v", key, err)
		}
		set[key] = mattr
	}
	return attributes.NewSet(set), nil
}
