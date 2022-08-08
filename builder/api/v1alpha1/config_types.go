package v1alpha1

// DataSetConfigurationKind object kind of DataSetConfiguration.
const DataSetConfigurationKind = "DataSetConfiguration"

// DataSetConfiguration configures a dataset
type DataSetConfiguration struct {
	Kind       string `mapstructure:"kind,omitempty"`
	APIVersion string `mapstructure:"apiVersion,omitempty"`
	// Collection configuration spec.
	Collection Collection `mapstructure:"collection,omitempty"`
	// LinkedCollections are the remote addresses of collection that are
	// linked to the collection.
	LinkedCollections []string `mapstructure:"linkedCollections,omitempty"`
}

type Collection struct {
	// Files defines custom attributes to add the the files in the
	// workspaces when publishing content/
	Files []File `mapstructure:"files,omitempty"`
	// SchemaAddress is the address of the schema to associated
	// to the Collection.
	SchemaAddress string `mapstructure:"file,omitempty"`
}

// File associates attributes with file names.
type File struct {
	// File is a string that can be compiled into a regular expression
	// for grouping attributes.
	File string `mapstructure:"file,omitempty"`
	// Attributes is the lists of to associate to the file.
	Attributes map[string]interface{} `mapstructure:"attributes,omitempty"`
}
