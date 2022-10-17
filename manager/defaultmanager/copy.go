package defaultmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	load "github.com/uor-framework/uor-client-go/config"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/ocimanifest"
	"github.com/uor-framework/uor-client-go/registryclient"
	"github.com/uor-framework/uor-client-go/util/workspace"
	"oras.land/oras-go/v2/content/memory"
)

func (d DefaultManager) Update(ctx context.Context, space workspace.Workspace, src string, dest string, add bool, remove bool, client registryclient.Client) (string, error) {

	// Pull the dataset config from the target collection
	_, manBytes, err := client.GetManifest(ctx, src)
	if err != nil {
		return "", err
	}

	var manifest ocispec.Manifest
	if err := json.NewDecoder(manBytes).Decode(&manifest); err != nil {
		return "", err
	}
	manconfig, err := client.GetContent(ctx, src, manifest.Config)
	if err != nil {
		return "", err
	}

	config, err := load.LoadDataSetConfig(manconfig)
	if err != nil {
		return "", err
	}

	if remove {
		pruneDescs, err := d.pullCollection(ctx, src, memory.New(), client)
		if err != nil {
			return "", err
		}
		pDesc := map[string]struct{}{}
		for _, s := range pruneDescs {
			pDesc[s.Digest.String()] = struct{}{}
		}

		var newLayers []ocispec.Descriptor
		for _, v := range manifest.Layers {
			if _, found := pDesc[v.Digest.String()]; !found {
				newLayers = append(newLayers, v)
			}
			manifest.Layers = newLayers
		}
	}
	var attributesByFile map[string]model.AttributeSet
	var files []string

	if add {
		files, attributesByFile, err = addFiles(space, d, config, dest)
		if err != nil {
			return "", err
		}
	}

	// If a schema is present, pull it and do the validation before
	// processing the files to get quick feedback to the user.
	collectionManifestAnnotations := map[string]string{}
	var descs []ocispec.Descriptor

	if config.Collection.SchemaAddress != "" {
		d.logger.Infof("Validating dataset configuration against schema %s", config.Collection.SchemaAddress)
		collectionManifestAnnotations[ocimanifest.AnnotationSchema] = config.Collection.SchemaAddress
		if err != nil {
			return "", fmt.Errorf("error configuring client: %v", err)
		}

		_, _, err = client.Pull(ctx, config.Collection.SchemaAddress, d.store)
		if err != nil {
			return "", err
		}

		schemaDoc, err := fetchJSONSchema(ctx, config.Collection.SchemaAddress, d.store)
		if err != nil {
			return "", err
		}
		if add {
			for file, attr := range attributesByFile {
				valid, err := schemaDoc.Validate(attr)
				if err != nil {
					return "", fmt.Errorf("schema validation error: %w", err)
				}
				if !valid {
					return "", fmt.Errorf("attributes for file %s are not valid for schema %s", file, config.Collection.SchemaAddress)
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

			descs, err = client.AddFiles(ctx, "", files...)
			if err != nil {
				return "", err
			}
		}

		descs = append(descs, manifest.Layers...)
		descs, err = ocimanifest.UpdateLayerDescriptors(descs, attributesByFile)
		if err != nil {
			return "", err
		}
	}

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

	linkedDescs, linkedSchemas, err := gatherLinkedCollections(ctx, config, client)
	if err != nil {
		return "", err
	}

	descs = append(descs, linkedDescs...)
	// Write the root collection attributes
	if len(linkedDescs) > 0 {
		collectionManifestAnnotations[ocimanifest.AnnotationSchemaLinks] = formatLinks(linkedSchemas)
		collectionManifestAnnotations[ocimanifest.AnnotationCollectionLinks] = formatLinks(config.Collection.LinkedCollections)
	}

	_, err = client.AddManifest(ctx, dest, configDesc, collectionManifestAnnotations, descs...)
	if err != nil {
		return "", err
	}

	desc, err := client.Save(ctx, dest, d.store)
	if err != nil {
		return "", fmt.Errorf("client save error for reference %s: %v", dest, err)
	}
	d.logger.Infof("Artifact %s built with reference name %s\n", desc.Digest, dest)

	return desc.Digest.String(), nil
}
