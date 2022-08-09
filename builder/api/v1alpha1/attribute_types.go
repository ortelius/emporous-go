package v1alpha1

// AttributeQueryKind object kind of AttributeQuery
const AttributeQueryKind = "AttributeQuery"

// AttributeQuery configures an attribute query.
type AttributeQuery struct {
	Kind       string `mapstructure:"kind"`
	APIVersion string `mapstructure:"apiVersion"`
	// Attributes list the configuration for Attribute types.
	Attributes Attributes `mapstructure:"attributes"`
}
