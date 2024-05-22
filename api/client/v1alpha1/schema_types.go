package v1alpha1

// SchemaConfigurationKind object kind of SchemaConfiguration
const SchemaConfigurationKind = "SchemaConfiguration"

// SchemaConfiguration configures a schema.
type SchemaConfiguration struct {
	TypeMeta `json:",inline"`
	Schema   SchemaConfigurationSpec `json:"schema"`
}

// SchemaConfigurationSpec defines the configuration spec to build an #mporous schema.
type SchemaConfigurationSpec struct {
	// ID is a name that will be used to identify
	// the schema
	ID          string `json:"id"`
	Description string `json:"description"`
	// SchemaPath defines that path to a JSON schema.
	SchemaPath string `json:"schemaPath"`
}
