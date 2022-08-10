package config

import (
	"fmt"
	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/builder/api/v1alpha1"
	"github.com/uor-framework/uor-client-go/model"
)

// ConvertToModel converts v1alpha1.Attributes to an model.AttributeSet.
func ConvertToModel(input v1alpha1.Attributes) (model.AttributeSet, error) {
	set := attributes.Attributes{}
	for key, val := range input {
		mattr, err := attributes.Reflect(key, val)
		if err != nil {
			return nil, fmt.Errorf("error converting attribute %s to model: %v", key, err)
		}
		set[key] = mattr
	}
	return set, nil
}
