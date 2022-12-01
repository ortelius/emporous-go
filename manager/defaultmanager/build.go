package defaultmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	empspec "github.com/emporous/collection-spec/specs-go/v1alpha1"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/registry"

	clientapi "github.com/emporous/emporous-go/api/client/v1alpha1"
	"github.com/emporous/emporous-go/attributes"
	"github.com/emporous/emporous-go/components"
	load "github.com/emporous/emporous-go/config"
	"github.com/emporous/emporous-go/content"
	"github.com/emporous/emporous-go/model"
	"github.com/emporous/emporous-go/nodes/descriptor"
	"github.com/emporous/emporous-go/nodes/descriptor/v2"
	"github.com/emporous/emporous-go/registryclient"
	"github.com/emporous/emporous-go/schema"
	"github.com/emporous/emporous-go/util/workspace"
)

// Build builds collection from input and store it in the underlying content store.
// If successful, the root descriptor is returned.
func (d DefaultManager) Build(ctx context.Context, space workspace.Workspace, config clientapi.DataSetConfiguration, reference string, client registryclient.Client) (string, error) {
	var files []string
	err := space.Walk(func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("traversing %s: %v", path, err)
		}
		if info == nil {
			return fmt.Errorf("no file info")
		}

		if info.Mode().IsRegular() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", fmt.Errorf("path %q empty workspace", space.Path("."))
	}

	var sets []model.AttributeSet
	regexpByFilename := map[string]*regexp.Regexp{}
	fileInfoByName := map[string]fileInformation{}
	for _, file := range config.Collection.Files {
		// Process each key into a regular expression and store it.

		// If the config has a grouping declared, make a valid regex.
		var expression string
		if strings.Contains(file.File, "*") && !strings.Contains(file.File, ".*") {
			expression = strings.Replace(file.File, "*", ".*", -1)
		} else {
			expression = strings.Replace(file.File, file.File, "^"+file.File+"$", -1)
		}

		nameSearch, err := regexp.Compile(expression)
		if err != nil {
			return "", err
		}
		regexpByFilename[file.File] = nameSearch

		set, err := load.ConvertToModel(file.Attributes)
		if err != nil {
			return "", err
		}
		sets = append(sets, set)

		fileInfo := fileInformation{
			AttributeSet: set,
			File:         file.FileInfo,
		}

		fileInfoByName[file.File] = fileInfo
	}

	// Merge the sets to ensure the dataset configuration
	// meet the schema require.
	mergedSet, err := attributes.Merge(sets...)
	if err != nil {
		return "", fmt.Errorf("failed to merge attributes: %w", err)
	}

	// If a schema is present, pull it and do the validation before
	// processing the files to get quick feedback to the user. Also, collection the schema ID
	// to place in the descriptor properties.
	schemaID := schema.UnknownSchemaID
	if config.Collection.SchemaAddress != "" {
		d.logger.Infof("Validating dataset configuration against schema %s", config.Collection.SchemaAddress)

		_, _, err := client.Pull(ctx, config.Collection.SchemaAddress, d.store)
		if err != nil {
			return "", fmt.Errorf("error configuring client: %v", err)
		}

		schemaDoc, detectedSchemaID, err := fetchJSONSchema(ctx, config.Collection.SchemaAddress, d.store)
		if err != nil {
			return "", err
		}

		if detectedSchemaID != "" {
			schemaID = detectedSchemaID
		}

		valid, err := schemaDoc.Validate(mergedSet)
		if err != nil {
			return "", fmt.Errorf("schema validation error: %w", err)
		}
		if !valid {
			return "", fmt.Errorf("attributes are not valid for schema %s", config.Collection.SchemaAddress)
		}
	}

	// To allow the files to be loaded relative to the render
	// workspace, change to the render directory. This is required
	// to get path correct in the description annotations.
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if err := os.Chdir(space.Path()); err != nil {
		return "", err
	}
	defer func() {
		if err := os.Chdir(cwd); err != nil {
			d.logger.Errorf("%v", err)
		}
	}()

	descs, err := client.AddFiles(ctx, "", files...)
	if err != nil {
		return "", err
	}

	// Gather workspace file metadata
	input := fmt.Sprintf("dir:%s", ".")
	inv, err := components.GenerateInventory(input, config)
	if err != nil {
		return "", fmt.Errorf("inventory generation for %s: %w", space.Path(), err)
	}

	// Create nodes and update node properties
	var nodes []v2.Node
	for _, desc := range descs {
		location, ok := desc.Annotations[ocispec.AnnotationTitle]
		if !ok {
			continue
		}
		// Using location as ID in this case because it is unique and
		// the digest may not be.
		node, err := v2.NewNode(location, desc)
		if err != nil {
			return "", err
		}
		node.Location = location
		if err := components.InventoryToProperties(*inv, location, node.Properties); err != nil {
			return "", err
		}
		nodes = append(nodes, *node)
	}

	updateFN := func(node v2.Node) error {
		if node.Location == "" {
			return nil
		}

		var sets []model.AttributeSet
		var fileConfig []empspec.File
		for file, fileInfo := range fileInfoByName {
			nameSearch := regexpByFilename[file]
			if nameSearch.Match([]byte(node.Location)) {
				if fileInfo.HasAttributes() {
					sets = append(sets, fileInfo.AttributeSet)
				}
				if fileInfo.HasFileInfo() {
					fileConfig = append(fileConfig, fileInfo.File)
				}
			}
		}

		switch {
		case len(fileConfig) == 1:
			node.Properties.File = &fileConfig[0]
		case len(fileConfig) > 1:
			return fmt.Errorf("file %q: more than one match for file configuration", node.Location)
		}

		merged, err := attributes.Merge(sets...)
		if err != nil {
			return err
		}
		if err := node.Properties.Merge(map[string]model.AttributeSet{schemaID: merged}); err != nil {
			return fmt.Errorf("file %s: %w", node.Location, err)
		}
		return nil
	}

	// Add user provided attributes to node properties
	descs, err = v2.UpdateDescriptors(nodes, updateFN)
	if err != nil {
		return "", err
	}

	// Store the DataSetConfiguration file in the manifest config of the OCI artifact for
	// later use.
	// Artifacts don't have configs. This will have to go with the regular descriptors.
	configJSON, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	configDesc, err := client.AddContent(ctx, empspec.MediaTypeConfiguration, configJSON, nil)
	if err != nil {
		return "", err
	}

	// Build index manifest
	manifestAnnotations := map[string]string{}
	if len(config.Collection.LinkedCollections) != 0 {
		aggregateDesc, err := d.addLinks(ctx, client, config.Collection.LinkedCollections)
		if err != nil {
			return "", err
		}
		aggregateDescJSON, err := json.Marshal(aggregateDesc)
		if err != nil {
			return "", err
		}
		manifestAnnotations[empspec.AnnotationLink] = string(aggregateDescJSON)
	}

	var prop descriptor.Properties
	// Add user specified component information to the manifest, if applicable.
	if config.Collection.Components.Name != "" {
		d.logger.Debugf("Component information detected. Adding inder core-descriptor schema.")
		componentAttr := &empspec.DescriptorAttributes{
			Component: empspec.Component{
				Name:      config.Collection.Components.Name,
				Version:   config.Collection.Components.Version,
				Type:      config.Collection.Components.Type,
				FoundBy:   config.Collection.Components.FoundBy,
				Locations: config.Collection.Components.Locations,
				Licenses:  config.Collection.Components.Licenses,
				Language:  config.Collection.Components.Language,
				CPEs:      config.Collection.Components.CPEs,
				PURL:      config.Collection.Components.PURL,
			},
		}
		prop.Descriptor = componentAttr
	}

	// Add user specified runtime information to the manifest, if applicable.
	if len(config.Collection.Runtime.Cmd) != 0 || len(config.Collection.Runtime.Entrypoint) != 0 {
		d.logger.Debugf("Runtime attributes detected. Adding under core-runtime schema")
		prop.Runtime = &config.Collection.Runtime
	}

	propsJSON, err := prop.MarshalJSON()
	if err != nil {
		return "", err
	}
	manifestAnnotations[empspec.AnnotationEmporousAttributes] = string(propsJSON)

	_, err = client.AddManifest(ctx, reference, configDesc, manifestAnnotations, descs...)
	if err != nil {
		return "", err
	}

	desc, err := client.Save(ctx, reference, d.store)
	if err != nil {
		return "", fmt.Errorf("client save error for reference %s: %v", reference, err)
	}
	d.logger.Infof("Artifact %s built with reference name %s\n", desc.Digest, reference)

	return desc.Digest.String(), nil
}

func (d DefaultManager) addLinks(ctx context.Context, client registryclient.Client, links []string) ([]ocispec.Descriptor, error) {
	d.logger.Infof("Processing %d link(s)", len(links))
	var linkedDesc []ocispec.Descriptor
	for _, l := range links {
		desc, _, err := client.GetManifest(ctx, l)
		if err != nil {
			return nil, fmt.Errorf("link %q: %w", l, err)
		}
		if desc.Annotations == nil {
			desc.Annotations = map[string]string{}
		}

		ref, err := registry.ParseReference(l)
		if err != nil {
			return nil, fmt.Errorf("link %q: %w", l, err)
		}
		linkAttr := descriptor.Properties{
			Link: &empspec.LinkAttributes{
				RegistryHint:  ref.Registry,
				NamespaceHint: ref.Repository,
				Transitive:    true,
			},
		}
		linkJSON, err := json.Marshal(linkAttr)
		if err != nil {
			return nil, err
		}
		desc.Annotations[empspec.AnnotationEmporousAttributes] = string(linkJSON)
		linkedDesc = append(linkedDesc, desc)
	}
	return linkedDesc, nil
}

// fetchJSONSchema returns a schema type from a content store and a schema address.
func fetchJSONSchema(ctx context.Context, schemaAddress string, store content.AttributeStore) (schema.Schema, string, error) {
	desc, err := store.AttributeSchema(ctx, schemaAddress)
	if err != nil {
		return schema.Schema{}, "", err
	}

	var schemaID string
	node, err := v2.NewNode(desc.Digest.String(), desc)
	if err != nil {
		return schema.Schema{}, "", err
	}
	props := node.Properties
	if props.IsASchema() {
		schemaID = props.Schema.ID
	}

	schemaReader, err := store.Fetch(ctx, desc)
	if err != nil {
		return schema.Schema{}, "", fmt.Errorf("error fetching schema from store: %w", err)
	}
	schemaBytes, err := ioutil.ReadAll(schemaReader)
	if err != nil {
		return schema.Schema{}, "", err
	}
	loader, err := schema.FromBytes(schemaBytes)
	if err != nil {
		return schema.Schema{}, "", err
	}

	sc, err := schema.New(loader)
	return sc, schemaID, err
}

// fileInformation pairs information configurable
// file attributes for comparison.
type fileInformation struct {
	model.AttributeSet
	empspec.File
}

func (f fileInformation) HasAttributes() bool {
	return f.AttributeSet.Len() > 0
}

func (f fileInformation) HasFileInfo() bool {
	return f.Permissions != 0 || f.UID != -1 || f.GID != -1
}
