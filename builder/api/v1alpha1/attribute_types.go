package v1alpha1

// AttributeQueryKind object kind of AttributeQuery
const AttributeQueryKind = "AttributeQuery"

// AttributeQuery configures an attribute query.
type AttributeQuery struct {
	TypeMeta `json:",inline"`
	// Attributes list the configuration for Attribute types.
	Attributes Attributes `json:"attributes"`
}
