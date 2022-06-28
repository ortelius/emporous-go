package attributes

import (
	"fmt"
	"sort"
	"strings"

	"github.com/uor-framework/client/model"
)

var (
	_ model.Matcher = &PartialAttributeMatcher{}
	_ model.Matcher = &ExactAttributeMatcher{}
)

// PartialAttributeMatcher contains configuration data for searching for a node by attribute.
// This matcher will check that the node attributes
type PartialAttributeMatcher map[string]string

// String list all attributes in the Matcher in a string format.
func (m PartialAttributeMatcher) String() string {
	return renderMatcher(m)
}

// Matches determines whether a node has all required attributes.
func (m PartialAttributeMatcher) Matches(n model.Node) bool {
	attr := n.Attributes()
	if attr == nil {
		return false
	}
	for key, value := range m {
		if exist := attr.Exists(key, value); !exist {
			return false
		}
	}
	return true
}

// ExactAttributeMatcher contains configuration data for searching for a node by attribute.
type ExactAttributeMatcher map[string]string

// String list all attributes in the Matcher in a string format.
func (m ExactAttributeMatcher) String() string {
	return renderMatcher(m)
}

// Matches determines whether a node has all required attributes.
func (m ExactAttributeMatcher) Matches(n model.Node) bool {
	attr := n.Attributes()
	if attr == nil {
		return false
	}
	if len(m) != attr.Len() {
		return false
	}
	for key, value := range m {
		if exist := attr.Exists(key, value); !exist {
			return false
		}
	}
	return true
}

// renderMatcher will render an attribute matcher as a string
func renderMatcher(m map[string]string) string {
	out := new(strings.Builder)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		line := fmt.Sprintf("%s=%s,", key, m[key])
		out.WriteString(line)
	}
	return strings.TrimSuffix(out.String(), ",")
}
