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
	"oras.land/oras-go/v2/content/file"
)

func (d DefaultManager) Update(ctx context.Context, space workspace.Workspace, dest string, add bool, remove bool, client registryclient.Client) (string, error) {

	// Pull the dataset config from the target collection
	_, manBytes, err := client.GetManifest(ctx, dest)
	if err != nil {
		return "", err
	}

	var manifest ocispec.Manifest
	if err := json.NewDecoder(manBytes).Decode(&manifest); err != nil {
		return "", err
	}
	manconfig, err := client.GetContent(ctx, dest, manifest.Config)
	if err != nil {
		return "", err
	}

	config, err := load.LoadDataSetConfig(manconfig)
	if err != nil {
		return "", err
	}
	d.logger.Infof("Dataset-config %s", config)
	attributesByFile := map[string]model.AttributeSet{}

	if remove {
		pruneDescs, err := d.pullCollection(ctx, dest, file.New("./output"), client)
		if err != nil {
			return "", err
		}
		d.logger.Infof("prune layers: %s", pruneDescs)
		for _, s := range pruneDescs {
			for i, v := range manifest.Layers {
				if v.Digest == s.Digest {
					d.logger.Infof("removing: %s", s.Digest)
					manifest.Layers = append(manifest.Layers[:i], manifest.Layers[i+1:]...)
				}
			}
		}
	}

	var files []string
	if add {
		err = space.Walk(func(path string, info os.FileInfo, err error) error {
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

		for _, file := range config.Collection.Files {
			set, err := load.ConvertToModel(file.Attributes)
			if err != nil {
				return "", err
			}
			attributesByFile[file.File] = set
		}

		if len(files) == 0 {
			return "", fmt.Errorf("path %q empty workspace", space.Path("."))
		}
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
	var descs []ocispec.Descriptor
	if add {
		descs, err = client.AddFiles(ctx, "", files...)
		if err != nil {
			return "", err
		}

		d.logger.Infof("Pre-update layers %s", manifest.Layers)
		descs = append(descs, manifest.Layers...)
		d.logger.Infof("Update layers %s", descs)
		descs, err = ocimanifest.UpdateLayerDescriptors(descs, attributesByFile)
		if err != nil {
			return "", err
		}
		d.logger.Infof("updated layer descs: %s", descs)
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
	d.logger.Infof("with links: %s", descs)
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
