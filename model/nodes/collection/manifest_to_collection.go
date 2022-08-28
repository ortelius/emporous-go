package collection

import (
	"context"
	"encoding/json"

	"github.com/google/go-containerregistry/pkg/v1/types"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	artifactspec "github.com/oras-project/artifacts-spec/specs-go/v1"

	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/model/nodes/descriptor"
	"github.com/uor-framework/uor-client-go/model/traversal"
)

// FetcherFunc fetches content for the specified descriptor
type FetcherFunc func(context.Context, ocispec.Descriptor) ([]byte, error)

// LoadFromManifest loads an OCI DAG into a Collection.
func LoadFromManifest(ctx context.Context, graph *Collection, fetcher FetcherFunc, manifest ocispec.Descriptor) error {
	// prepare pre-handler
	root, err := descriptor.NewNode(manifest.Digest.String(), manifest)
	if err != nil {
		return err
	}

	// track content status
	tracker := traversal.NewTracker(root, nil)

	handler := traversal.HandlerFunc(func(ctx context.Context, tracker traversal.Tracker, node model.Node) ([]model.Node, error) {
		// skip the node if it has been indexed
		if graph.HasNode(node.ID()) {
			return nil, traversal.ErrSkip
		}

		desc, ok := node.(*descriptor.Node)
		if !ok {
			return nil, traversal.ErrSkip
		}

		successors, err := getSuccessors(ctx, fetcher, manifest)
		if err != nil {
			return nil, err
		}

		nodes, err := indexNode(graph, desc.Descriptor(), successors)
		if err != nil {
			return nil, err
		}

		return nodes, nil
	})

	return tracker.Walk(ctx, handler, root)
}

// AddManifest will add a single manifest to the Collection.
func AddManifest(ctx context.Context, graph *Collection, fetcher FetcherFunc, node ocispec.Descriptor) error {
	successors, err := getSuccessors(ctx, fetcher, node)
	if err != nil {
		return err
	}
	if _, err := indexNode(graph, node, successors); err != nil {
		return err
	}
	return nil
}

// indexNode indexes relationships between child and parent nodes.
func indexNode(graph *Collection, node ocispec.Descriptor, successors []ocispec.Descriptor) ([]model.Node, error) {
	n, err := addOrGetNode(graph, node)
	if err != nil {
		return nil, err
	}
	var result []model.Node
	for _, successor := range successors {
		s, err := addOrGetNode(graph, successor)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
		e := NewEdge(n, s)
		if err := graph.AddEdge(e); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// addOrGetNode will return the node if it exists in the graph or will create a new
// descriptor node.
func addOrGetNode(graph *Collection, desc ocispec.Descriptor) (model.Node, error) {
	n := graph.NodeByID(desc.Digest.String())
	if n != nil {
		return n, nil
	}
	n, err := descriptor.NewNode(desc.Digest.String(), desc)
	if err != nil {
		return n, err
	}
	if err := graph.AddNode(n); err != nil {
		return nil, err
	}
	return n, nil
}

// getSuccessor returns the nodes directly pointed by the current node. This is adapted from the `oras` content.Successors
// to allow the use of a function signature to pull descriptor content.
func getSuccessors(ctx context.Context, fetcher FetcherFunc, node ocispec.Descriptor) ([]ocispec.Descriptor, error) {
	switch node.MediaType {
	case string(types.DockerManifestSchema2), ocispec.MediaTypeImageManifest:
		content, err := fetcher(ctx, node)
		if err != nil {
			return nil, err
		}

		// docker manifest and oci manifest are equivalent for successors.
		var manifest ocispec.Manifest
		if err := json.Unmarshal(content, &manifest); err != nil {
			return nil, err
		}
		return append([]ocispec.Descriptor{manifest.Config}, manifest.Layers...), nil
	case string(types.DockerManifestList), ocispec.MediaTypeImageIndex:
		content, err := fetcher(ctx, node)
		if err != nil {
			return nil, err
		}

		// docker manifest list and oci index are equivalent for successors.
		var index ocispec.Index
		if err := json.Unmarshal(content, &index); err != nil {
			return nil, err
		}

		return index.Manifests, nil
	case artifactspec.MediaTypeArtifactManifest:
		content, err := fetcher(ctx, node)
		if err != nil {
			return nil, err
		}

		var manifest artifactspec.Manifest
		if err := json.Unmarshal(content, &manifest); err != nil {
			return nil, err
		}
		var nodes []ocispec.Descriptor
		if manifest.Subject != nil {
			nodes = append(nodes, artifactToOCI(*manifest.Subject))
		}
		for _, blob := range manifest.Blobs {
			nodes = append(nodes, artifactToOCI(blob))
		}
		return nodes, nil
	}
	return nil, nil
}

// artifactToOCI converts artifact descriptor to OCI descriptor.
func artifactToOCI(desc artifactspec.Descriptor) ocispec.Descriptor {
	return ocispec.Descriptor{
		MediaType:   desc.MediaType,
		Digest:      desc.Digest,
		Size:        desc.Size,
		URLs:        desc.URLs,
		Annotations: desc.Annotations,
	}
}
