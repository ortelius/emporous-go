package defaultmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	clientapi "github.com/uor-framework/uor-client-go/api/client/v1alpha1"
	load "github.com/uor-framework/uor-client-go/config"
	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/ocimanifest"
	"github.com/uor-framework/uor-client-go/registryclient"
	"github.com/uor-framework/uor-client-go/schema"
	"github.com/uor-framework/uor-client-go/util/workspace"
)

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

	attributesByFile := map[string]model.AttributeSet{}
	for _, file := range config.Collection.Files {
		set, err := load.ConvertToModel(file.Attributes)
		if err != nil {
			return "", err
		}
		attributesByFile[file.File] = set
	}

	// If a schema is present, pull it and do the validation before
	// processing the files to get quick feedback to the user.
	collectionManifestAnnotations := map[string]string{}
	if config.Collection.SchemaAddress != "" {
		d.logger.Infof("Validating dataset configuration against schema %s", config.Collection.SchemaAddress)
		collectionManifestAnnotations[ocimanifest.AnnotationSchema] = config.Collection.SchemaAddress
		if err != nil {
			return "", fmt.Errorf("error configuring client: %v", err)
		}

		_, err = client.Pull(ctx, config.Collection.SchemaAddress, d.store)
		if err != nil {
			return "", err
		}

		schemaDoc, err := fetchJSONSchema(ctx, config.Collection.SchemaAddress, d.store)
		if err != nil {
			return "", err
		}

		for file, attr := range attributesByFile {
			valid, err := schemaDoc.Validate(attr)
			if err != nil {
				return "", fmt.Errorf("schema validation error: %w", err)
			}
			if !valid {
				return "", fmt.Errorf("attributes for file %s are not valid for schema %s", file, config.Collection.SchemaAddress)
			}
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

	descs, err = ocimanifest.UpdateLayerDescriptors(descs, attributesByFile)
	if err != nil {
		return "", err
	}

	linkedDescs, linkedSchemas, err := gatherLinkedCollections(ctx, config, client)
	if err != nil {
		return "", err
	}

	descs = append(descs, linkedDescs...)

	// Store the DataSetConfiguration file in the manifest config of the OCI artifact for
	// later use.
	configJSON, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	configDesc, err := client.AddContent(ctx, ocimanifest.UORConfigMediaType, configJSON, nil)
	if err != nil {
		return "", err
	}

	// Write the root collection attributes
	if len(linkedDescs) > 0 {
		collectionManifestAnnotations[ocimanifest.AnnotationSchemaLinks] = formatLinks(linkedSchemas)
		collectionManifestAnnotations[ocimanifest.AnnotationCollectionLinks] = formatLinks(config.Collection.LinkedCollections)
	}

	_, err = client.AddManifest(ctx, reference, configDesc, collectionManifestAnnotations, descs...)
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

// fetchJSONSchema returns a schema type from a content store and a schema address.
func fetchJSONSchema(ctx context.Context, schemaAddress string, store content.AttributeStore) (schema.Schema, error) {
	desc, err := store.AttributeSchema(ctx, schemaAddress)
	if err != nil {
		return schema.Schema{}, err
	}
	schemaReader, err := store.Fetch(ctx, desc)
	if err != nil {
		return schema.Schema{}, fmt.Errorf("error fetching schema from store: %w", err)
	}
	schemaBytes, err := ioutil.ReadAll(schemaReader)
	if err != nil {
		return schema.Schema{}, err
	}
	return schema.FromBytes(schemaBytes)
}

// gatherLinkedCollections create null descriptors to denotes linked collections in a manifest with schema link information.
func gatherLinkedCollections(ctx context.Context, cfg clientapi.DataSetConfiguration, client registryclient.Client) ([]ocispec.Descriptor, []string, error) {
	var allLinkedSchemas []string
	var linkedDescs []ocispec.Descriptor
	for _, collection := range cfg.Collection.LinkedCollections {

		rootSchema, linkedSchemas, err := getLinks(ctx, collection, client)
		if err != nil {
			return nil, nil, fmt.Errorf("collection %q: %w", collection, err)
		}

		if len(linkedSchemas) != 0 {
			allLinkedSchemas = append(allLinkedSchemas, linkedSchemas...)
		}

		allLinkedSchemas = append(allLinkedSchemas, rootSchema)

		annotations := map[string]string{
			ocimanifest.AnnotationSchema:      rootSchema,
			ocimanifest.AnnotationSchemaLinks: formatLinks(linkedSchemas),
		}
		// The bytes contain the collection name to keep the blobs unique within the manifest
		desc, err := client.AddContent(ctx, ocispec.MediaTypeImageLayer, []byte(collection), annotations)
		if err != nil {
			return nil, nil, err
		}
		linkedDescs = append(linkedDescs, desc)
	}
	return linkedDescs, allLinkedSchemas, nil
}

// getLinks retrieves all schema information for a given reference.
func getLinks(ctx context.Context, reference string, client registryclient.Remote) (string, []string, error) {
	_, manBytes, err := client.GetManifest(ctx, reference)
	if err != nil {
		return "", nil, err
	}
	defer manBytes.Close()
	return ocimanifest.FetchSchemaLinks(manBytes)
}

func formatLinks(links []string) string {
	n := len(links)
	switch {
	case n == 1:
		return links[0]
	case n > 1:
		dedupLinks := deduplicate(links)
		return strings.Join(dedupLinks, ocimanifest.Separator)
	default:
		return ""
	}
}

func deduplicate(in []string) []string {
	links := map[string]struct{}{}
	var out []string
	for _, l := range in {
		if _, ok := links[l]; ok {
			continue
		}
		links[l] = struct{}{}
		out = append(out, l)
	}
	return out
}
