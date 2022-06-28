package v1alpha1

// DataSetConfigurationKind object kind of DataSetConfiguration.
const DataSetConfigurationKind = "DataSetConfiguration"

// DataSetConfiguration configures a dataset
type DataSetConfiguration struct {
	Kind       string `mapstructure:"kind,omitempty"`
	APIVersion string `mapstructure:"apiVersion,omitempty"`
	Files      []File `mapstructure:"files,omitempty"`
}

// File associates attributes with file names
type File struct {
	File       string            `mapstructure:"file,omitempty"`
	Attributes map[string]string `mapstructure:"attributes,omitempty"`
}
