package v1alpha1

// AttributeQueryKind object kind of AttributeQuery
const AttributeQueryKind = "AttributeQuery"

// AttributeQuery configures an attribute query.
type AttributeQuery struct {
	Kind       string `mapstructure:"kind"`
	APIVersion string `mapstructure:"apiVersion"`
	// Attributes list the configuration for Attribute types.
	Attributes []Attribute `mapstructure:"attributes"`
}

// Attribute construct a query for an individual attribute.
type Attribute struct {
	// Key represent the attribute key.
	Key string `mapstructure:"key"`
	// Value represent an attribute value.
	Value interface{} `mapstructure:"value"`
}
