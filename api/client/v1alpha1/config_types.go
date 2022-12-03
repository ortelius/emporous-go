package v1alpha1

import (
	"encoding/json"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"
)

// DataSetConfigurationKind object kind of DataSetConfiguration.
const DataSetConfigurationKind = "DataSetConfiguration"

// DataSetConfiguration configures a dataset
type DataSetConfiguration struct {
	TypeMeta `json:",inline"`
	// Collection configuration spec.
	Collection DataSetConfigurationSpec `json:"collection,omitempty"`
}

// DataSetConfigurationSpec defines the configuration spec to build a single UOR collection.
type DataSetConfigurationSpec struct {
	// Components attaches component information to specific files.
	Components ComponentSpec `json:"components,omitempty"`
	// Runtime attaches runtime information to the artifact manifest
	Runtime ocispec.ImageConfig `json:"runtime,omitempty"`
	// Files defines custom attributes to add the files in the
	// workspaces when publishing content/
	Files []File `json:"files,omitempty"`
	// SchemaAddress is the address of the schema to associated
	// to the Collection.
	SchemaAddress string `json:"schemaAddress,omitempty"`
	// LinkedCollections are the remote addresses of collection that are
	// linked to the collection.
	LinkedCollections []string `json:"linkedCollections,omitempty"`
}

// ComponentSpec defines configuration information when creating component lists.
// Each field except for Platform is will allow users to set manifest-level component
// information. All workspace items will have their component information collection on a best-effort
// basis.
type ComponentSpec struct {
	Platform  string   `json:"platform"`
	Name      string   `json:"name"`
	Version   string   `json:"version"`
	Type      string   `json:"type"`
	FoundBy   string   `json:"foundBy"`
	Locations []string `json:"locations"`
	Licenses  []string `json:"licenses"`
	Language  string   `json:"language"`
	// Common Platform Enumeration
	CPEs []string `json:"cpes"`
	// Package URL
	PURL string `json:"purl"`
}

// File associates attributes with file names.
type File struct {
	// File is a string that can be compiled into a regular expression
	// for grouping attributes.
	File string `json:"file,omitempty"`
	// FileInfo sets target path, ownership, and
	// permissions for files that can be used with container runtimes.
	FileInfo uorspec.File `json:"fileInfo,omitempty"`
	// Attributes is the lists of to associate to the file.
	Attributes Attributes `json:"attributes,omitempty"`
}

// Attributes is a map structure that holds all
// attribute information provided by the user.
type Attributes map[string]interface{}

// UnmarshalJSON sets custom unmarshalling logic to File.
// In this case it sets the default UID and GID to invalid
// ID numbers to differentiate between values intentionally set at 0.
func (f *File) UnmarshalJSON(data []byte) error {
	type fileAlias File
	test := &fileAlias{
		FileInfo: uorspec.File{
			UID: -1,
			GID: -1,
		},
	}

	err := json.Unmarshal(data, test)
	if err != nil {
		return err
	}
	*f = File(*test)
	return nil
}
