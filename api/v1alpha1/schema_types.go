package v1alpha1

import (
	"github.com/uor-framework/uor-client-go/schema"
)

// SchemaConfigurationKind object kind of SchemaConfiguration
const SchemaConfigurationKind = "SchemaConfiguration"

// SchemaConfiguration configures a schema.
type SchemaConfiguration struct {
	TypeMeta `json:",inline"`
	// Address is the remote location for the default schema of the
	// collection.
	Address string `json:"address"`
	// DefaultContentDeclarations defined that default arguments that the
	// Algorithm will use for processing.
	DefaultContentDeclarations map[string]string `json:"defaultContentDeclarations,omitempty"`
	// CommonAttributeMapping defines common attribute keys and values for schema. The values
	// must be in JSON Format.
	CommonAttributeMapping Attributes `json:"commonAttributeMapping,omitempty"`
	// AttributeTypes is a collection of attribute type definitions.
	AttributeTypes schema.Types `json:"attributeTypes,omitempty"`
}
