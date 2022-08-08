package matchers

import (
	"github.com/uor-framework/uor-client-go/model"
)

var (
	_ model.Matcher = &PartialAttributeMatcher{}
)

// PartialAttributeMatcher contains configuration data for searching for a node by attribute.
// This matcher will check that the node attributes
type PartialAttributeMatcher map[string]model.Attribute

// Matches determines whether a node has all required attributes.
func (m PartialAttributeMatcher) Matches(n model.Node) bool {
	attr := n.Attributes()
	if attr == nil {
		return false
	}

	for key, value := range m {
		if exist := attr.Exists(key, value.Kind(), value.AsAny()); !exist {
			return false
		}
	}
	return true
}
