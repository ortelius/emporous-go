package ocimanifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/uor-framework/uor-client-go/api/v1alpha1"
	"github.com/uor-framework/uor-client-go/config"
	"regexp"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
)

const (
	// AnnotationSchema is the reference to the
	// default schema of the collection.
	AnnotationSchema = "uor.schema"
	// AnnotationSchemaLinks is the reference to linked
	// schemas for a collection. This will define all referenced
	// schemas for the collection and sub-collection. The tree will
	// be fully resolved.
	AnnotationSchemaLinks = "uor.schema.linked"
	// AnnotationCollectionLinks references the collections
	// that are linked to a collection node. They will only
	// reference adjacent collection and will not descend
	// into sub-collections.
	AnnotationCollectionLinks = "uor.collections.linked"
	// AnnotationUORAttributes references the collection attributes in a
	// JSON format.
	AnnotationUORAttributes = "uor.attributes"
	// Separator is the value used to denote a list of
	// schema or collection in a manifest.
	Separator = ","
	// UORConfigMediaType is the manifest config media type
	// for UOR OCI manifests.
	UORConfigMediaType = "application/vnd.uor.config.v1+json"
)

var (
	// ErrNoKnownSchema denotes that no schema
	// annotation is set on the manifest.
	ErrNoKnownSchema = errors.New("no schema")
	// ErrNoCollectionLinks denotes that the manifest
	// does contain annotation that set collection links.
	ErrNoCollectionLinks = errors.New("no collection links")
)

// AnnotationsToAttributeSet converts annotations from descriptors
// to an AttributeSet. This also perform annotation validation.
func AnnotationsToAttributeSet(annotations map[string]string, skip func(string) bool) (model.AttributeSet, error) {
	set := attributes.Attributes{}

	for key, value := range annotations {
		if skip != nil && skip(key) {
			continue
		}

		// Key collision.
		// TODO(jpower432): Handle this more gracefully
		if _, exists := set[key]; exists {
			continue
		}

		if key != AnnotationUORAttributes {
			set[key] = attributes.NewString(key, value)
			continue
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(value), &data); err != nil {
			return set, err
		}
		for jKey, jVal := range data {
			attr, err := attributes.Reflect(jKey, jVal)
			if err != nil {
				return set, fmt.Errorf("annotation %q: error creating attribute: %w", key, err)
			}
			set[jKey] = attr
		}
	}
	return set, nil
}

// AnnotationsFromAttributeSet converts an AttributeSet to annotations. All annotation values
// are saved in a JSON valid syntax to allow for typing upon retrieval.
func AnnotationsFromAttributeSet(set model.AttributeSet) (map[string]string, error) {
	return map[string]string{AnnotationUORAttributes: string(set.AsJSON())}, nil
}

// UpdateLayerDescriptors updates layers descriptor annotations with user provided key,value pairs
func UpdateLayerDescriptors(descs []ocispec.Descriptor, cfg v1alpha1.DataSetConfiguration) ([]ocispec.Descriptor, error) {
	var updateDescs []ocispec.Descriptor
	for _, desc := range descs {
		filename, ok := desc.Annotations[ocispec.AnnotationTitle]
		if !ok {
			// skip any descriptor with no name attached
			continue
		}

		for _, file := range cfg.Collection.Files {
			// If the config has a grouping declared, make a valid regex.
			if strings.Contains(file.File, "*") && !strings.Contains(file.File, ".*") {
				file.File = strings.Replace(file.File, "*", ".*", -1)
			} else {
				file.File = strings.Replace(file.File, file.File, "^"+file.File+"$", -1)
			}
			namesearch, err := regexp.Compile(file.File)
			if err != nil {
				return nil, err
			}

			if namesearch.Match([]byte(filename)) {
				// Get the k/v pairs from the config and add them to the descriptor annotations.
				set, err := config.ConvertToModel(file.Attributes)
				if err != nil {
					return nil, err
				}
				annotations, err := AnnotationsFromAttributeSet(set)
				if err != nil {
					return nil, err
				}
				for key, value := range annotations {
					desc.Annotations[key] = value
				}
			}
		}
		updateDescs = append(updateDescs, desc)
	}

	return updateDescs, nil
}
