package descriptor

import (
	"encoding/json"
	"errors"

	"github.com/nsf/jsondiff"

	"github.com/emporous/emporous-go/model"
)

var _ model.Matcher = JSONSubsetMatcher{}

// JSONSubsetMatcher check that the node attributes are a superset or an
// exact match to the given json input. The JSONSubsetMatcher should be used when
// to filter descriptor type nodes to include core schema fields.
type JSONSubsetMatcher json.RawMessage

// Matches determines whether a node has all required attributes.
func (m JSONSubsetMatcher) Matches(n model.Node) (bool, error) {
	attr := n.Attributes()
	if attr == nil {
		return false, errors.New("node attributes cannot be nil")
	}

	attrJSON, err := attr.MarshalJSON()
	if err != nil {
		return false, err
	}

	opts := jsondiff.DefaultJSONOptions()

	res, _ := jsondiff.Compare(attrJSON, m, &opts)

	if res != jsondiff.NoMatch {
		return true, nil
	}

	return false, nil
}
