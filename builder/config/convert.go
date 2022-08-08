package config

import (
	"fmt"
	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/builder/api/v1alpha1"
	"github.com/uor-framework/uor-client-go/model"
)

// ConvertToModel converts an attribute query to an attribute set.
func ConvertToModel(query v1alpha1.AttributeQuery) (model.AttributeSet, error) {
	set := attributes.Attributes{}
	for _, attr := range query.Attributes {
		mattr, err := attr.ToModel()
		if err != nil {
			return nil, fmt.Errorf("error converting attribute %s to model", attr.Key)
		}
		set[attr.Key] = mattr
	}
	return set, nil
}
