package attributes

import (
	"fmt"
	"sort"
	"strings"

	"github.com/uor-framework/client/model"
)

// Attributes implements the model.Attributes interface
// using a multi-map storing a set of values.
// The current implementation would allow for aggregation of the attributes
// of child nodes to the parent nodes.
// TODO(jpower432): Research alternative data structures for storing these values.
// Since this will most likely limit the collection size
// there could we a more efficient data structure for storing multiple values
// for one key for data aggregation.
type Attributes map[string]map[string]struct{}

var _ model.Attributes = &Attributes{}

// Find returns all values stored for a specified key.
func (a Attributes) Find(key string) []string {
	valSet, exists := a[key]
	if !exists {
		return nil
	}
	var vals []string
	for val := range valSet {
		vals = append(vals, val)
	}
	return vals
}

// Exists returns whether a key,value pair exists in the
// attribute set.
func (a Attributes) Exists(key, value string) bool {
	vals, exists := a[key]
	if !exists {
		return false
	}
	_, valExists := vals[value]
	return valExists
}

// Strings returns a string representation of the
// attribute set.
func (a Attributes) String() string {
	out := new(strings.Builder)
	keys := make([]string, 0, len(a))
	for k := range a {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		vals := a.List()[key]
		sort.Strings(vals)
		for _, val := range vals {
			line := fmt.Sprintf("%s=%s,", key, val)
			out.WriteString(line)
		}
	}
	return strings.TrimSuffix(out.String(), ",")
}

// List will list all key, value pairs for the attributes in a
// consumable format.
func (a Attributes) List() map[string][]string {
	list := make(map[string][]string, len(a))
	for key, vals := range a {
		for val := range vals {
			list[key] = append(list[key], val)
		}
	}
	return list
}

// Len returns the length of the attribute set.
func (a Attributes) Len() int {
	return len(a)
}

// Merge will merge the input Attributes with the receiver.
func (a Attributes) Merge(attr model.Attributes) {
	for key, vals := range attr.List() {
		for _, val := range vals {
			sub := a[key]
			sub[val] = struct{}{}
		}
	}
}
