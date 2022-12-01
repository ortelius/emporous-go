package components

import (
	"fmt"

	"github.com/anchore/syft/cmd/syft/cli/eventloop"
	"github.com/anchore/syft/syft/artifact"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"

	clientapi "github.com/emporous/emporous-go/api/client/v1alpha1"
	"github.com/emporous/emporous-go/version"
)

const ApplicationName = "uor"

// GenerateInventory generates an inventory based on input and DatasetConfiguration information.
func GenerateInventory(input string, config clientapi.DataSetConfiguration) (*sbom.SBOM, error) {
	si, err := source.ParseInput(input, config.Collection.Components.Platform, true)
	if err != nil {
		return nil, fmt.Errorf("could not generate source input:  %w", err)
	}

	src, cleanup, err := source.New(*si, nil, nil)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	s, err := generateSBOM(src, config)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func generateSBOM(src *source.Source, config clientapi.DataSetConfiguration) (*sbom.SBOM, error) {
	tasks, err := getTasks(config)
	if err != nil {
		return nil, err
	}

	s := sbom.SBOM{
		Source: src.Metadata,
		Descriptor: sbom.Descriptor{
			Name:          ApplicationName,
			Version:       version.GetVersion(),
			Configuration: config,
		},
	}

	if err := buildRelationships(&s, src, tasks); err != nil {
		return nil, err
	}

	return &s, nil
}

func buildRelationships(s *sbom.SBOM, src *source.Source, tasks []eventloop.Task) error {
	var relationships []artifact.Relationship
	for _, task := range tasks {
		relationship, err := task(&s.Artifacts, src)
		if err != nil {
			return err
		}
		relationships = append(relationships, relationship...)
	}

	s.Relationships = append(s.Relationships, relationships...)
	return nil
}
