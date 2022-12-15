package matchers

import (
	"errors"
	"fmt"

	"github.com/emporous/emporous-go/model"
)

var _ model.Matcher = PartialAttributeMatcher{}

// PartialAttributeMatcher contains configuration data for searching for a node by attribute.
// This matcher will check that the node attributes
type PartialAttributeMatcher map[string]model.AttributeValue

// Matches determines whether a node has all required attributes.
func (m PartialAttributeMatcher) Matches(n model.Node) (bool, error) {
	attr := n.Attributes()
	if attr == nil {
		return false, errors.New("node attributes cannot be nil")
	}

	for key, a := range m {
		exist, err := attr.Exists(key, a)
		if err != nil {
			return false, fmt.Errorf("error evaluating attribute %s: %w", key, err)
		}
		if !exist {
			return false, nil
		}
	}
	return true, nil
}
