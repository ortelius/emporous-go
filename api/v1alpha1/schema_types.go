package v1alpha1

import (
	"github.com/uor-framework/uor-client-go/schema"
)

// SchemaConfigurationKind object kind of SchemaConfiguration
const SchemaConfigurationKind = "SchemaConfiguration"

// SchemaConfiguration configures a schema.
type SchemaConfiguration struct {
	TypeMeta `json:",inline"`
	Schema   SchemaConfigurationSpec `json:"schema"`
}

// SchemaConfigurationSpec defines the configuration spec to build a UOR schema.
type SchemaConfigurationSpec struct {
	// Address is the remote location for the default schema of the
	// collection.
	Address string `json:"address"`
	// AttributeTypes is a collection of attribute type definitions.
	AttributeTypes schema.Types `json:"attributeTypes,omitempty"`
}
