package matchers

import (
	"errors"
	"fmt"

	"github.com/uor-framework/uor-client-go/model"
)

var _ model.Matcher = PartialAttributeMatcher{}

// PartialAttributeMatcher contains configuration data for searching for a node by attribute.
// This matcher will check that the node attributes
type PartialAttributeMatcher map[string]model.Attribute

// Matches determines whether a node has all required attributes.
func (m PartialAttributeMatcher) Matches(n model.Node) (bool, error) {
	attr := n.Attributes()
	if attr == nil {
		return false, errors.New("node attributes cannot be nil")
	}

	for _, a := range m {
		exist, err := attr.Exists(a)
		if err != nil {
			return false, fmt.Errorf("error evaluating attribute %s: %w", a.Key(), err)
		}
		if !exist {
			return false, nil
		}
	}
	return true, nil
}
