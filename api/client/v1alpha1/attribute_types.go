package v1alpha1

import "encoding/json"

// AttributeQueryKind object kind of AttributeQuery
const AttributeQueryKind = "AttributeQuery"

// AttributeQuery configures an attribute query against
// emporous collection content.
type AttributeQuery struct {
	TypeMeta `json:",inline"`
	// Attributes list the configuration for Attribute types.
	Attributes json.RawMessage `json:"attributes"`
}
