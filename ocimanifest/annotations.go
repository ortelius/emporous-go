package ocimanifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/uor-framework/uor-client-go/api/v1alpha1"
	"github.com/uor-framework/uor-client-go/config"

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
	// UORSchemaMediaType is the media type for a UOR schema.
	UORSchemaMediaType = "application/vnd.uor.schema.v1+json"
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
// to an AttributeSet. This also performs annotation validation.
func AnnotationsToAttributeSet(annotations map[string]string, skip func(string) bool) (model.AttributeSet, error) {
	set := attributes.Attributes{}

	for key, value := range annotations {
		if skip != nil && skip(key) {
			continue
		}

		// Handle key collision. This should only occur if
		// an annotation is set and the key also exists in the UOR
		// specific attributes.
		// TODO(jpower432): Handle more gracefully.
		if _, exists := set[key]; exists {
			continue
		}

		// Since annotations are in the form of map[string]string, we
		// can just assume it is a string attribute at this point. Incorporating
		// this into thr attribute set allows, users to pull by filename or reference name (cache).
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

// UpdateLayerDescriptors updates layers descriptor annotations with attributes from an AttributeSet. The key in the fileAttributes
// argument can be a regular expression or the name of a single file.
func UpdateLayerDescriptors(descs []ocispec.Descriptor, fileAttributes map[string]model.AttributeSet) ([]ocispec.Descriptor, error) {
	// Process each key into a regular expression and store it.
	regexpByFilename := map[string]*regexp.Regexp{}
	for file := range fileAttributes {
		// If the config has a grouping declared, make a valid regex.
		var expression string
		if strings.Contains(file, "*") && !strings.Contains(file, ".*") {
			expression = strings.Replace(file, "*", ".*", -1)
		} else {
			expression = strings.Replace(file, file, "^"+file+"$", -1)
		}

		nameSearch, err := regexp.Compile(expression)
		if err != nil {
			return nil, err
		}
		regexpByFilename[file] = nameSearch
	}

	var updateDescs []ocispec.Descriptor
	for _, desc := range descs {
		filename, ok := desc.Annotations[ocispec.AnnotationTitle]
		if !ok {
			// skip any descriptor with no name attached
			continue
		}

		for file, set := range fileAttributes {
			nameSearch := regexpByFilename[file]
			if nameSearch.Match([]byte(filename)) {
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
