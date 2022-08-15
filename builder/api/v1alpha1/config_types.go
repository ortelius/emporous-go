package v1alpha1

// TypeMeta contains type metadata.
type TypeMeta struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

// DataSetConfigurationKind object kind of DataSetConfiguration.
const DataSetConfigurationKind = "DataSetConfiguration"

// DataSetConfiguration configures a dataset
type DataSetConfiguration struct {
	TypeMeta `json:",inline"`
	// Collection configuration spec.
	Collection Collection `json:"collection,omitempty"`
	// LinkedCollections are the remote addresses of collection that are
	// linked to the collection.
	LinkedCollections []string `json:"linkedCollections,omitempty"`
}

type Collection struct {
	// Files defines custom attributes to add the files in the
	// workspaces when publishing content/
	Files []File `json:"files,omitempty"`
	// SchemaAddress is the address of the schema to associated
	// to the Collection.
	SchemaAddress string `json:"schemaAddress,omitempty"`
}

// File associates attributes with file names.
type File struct {
	// File is a string that can be compiled into a regular expression
	// for grouping attributes.
	File string `json:"file,omitempty"`
	// Attributes is the lists of to associate to the file.
	Attributes Attributes `json:"attributes,omitempty"`
}

// Attributes is a map structure that holds all
// attribute information provided by the user.
type Attributes map[string]interface{}
