package defaultmanager

import (
	"context"
	"errors"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/model/traversal"
	"github.com/uor-framework/uor-client-go/ocimanifest"
	"github.com/uor-framework/uor-client-go/registryclient"
)

func (d DefaultManager) Pull(ctx context.Context, source string, remote registryclient.Remote, destination content.Store) ([]string, error) {
	descs, err := d.pullCollection(ctx, source, destination, remote)
	if err != nil {
		return nil, err
	}

	var digests []string
	for _, desc := range descs {
		digests = append(digests, desc.Digest.String())
		d.logger.Infof("Found %s", desc.Digest)
	}
	return digests, nil
}

func (d DefaultManager) PullAll(ctx context.Context, source string, remote registryclient.Remote, destination content.Store) ([]string, error) {
	root, err := remote.LoadCollection(ctx, source)
	if err != nil {
		return nil, err
	}
	descs, err := d.copyCollections(ctx, &root, destination, remote)
	if err != nil {
		return nil, err
	}

	var digests []string
	for _, desc := range descs {
		digests = append(digests, desc.Digest.String())
		d.logger.Infof("Found %s", desc.Digest)
	}
	return digests, nil
}

// pullCollection pulls a single collection and returns the manifest descriptors and an error.
func (d DefaultManager) pullCollection(ctx context.Context, reference string, destination content.Store, remote registryclient.Remote) ([]ocispec.Descriptor, error) {
	rootDesc, descs, err := remote.Pull(ctx, reference, destination)
	if err != nil {
		return nil, err
	}

	// Ensure the store is tagged with the new reference.
	if len(rootDesc.Digest) != 0 {
		return descs, d.store.Tag(ctx, rootDesc, reference)
	}

	return descs, nil
}

// copy performs graph traversal of linked collections and performs collection copies filtered by the matcher.
func (d DefaultManager) copyCollections(ctx context.Context, root model.Node, destination content.Store, remote registryclient.Remote) ([]ocispec.Descriptor, error) {
	seen := map[string]struct{}{}
	var allDescs []ocispec.Descriptor

	tracker := traversal.NewTracker(root, nil)
	handler := traversal.HandlerFunc(func(ctx context.Context, tracker traversal.Tracker, node model.Node) ([]model.Node, error) {

		descs, err := d.pullCollection(ctx, node.Address(), destination, remote)
		if err != nil {
			return nil, err
		}
		allDescs = append(allDescs, descs...)

		successors, err := getSuccessors(ctx, node.Address(), remote)
		if err != nil {
			if errors.Is(err, ocimanifest.ErrNoCollectionLinks) {
				d.logger.Debugf("collection %s has no links", node.Address())
				return nil, nil
			}
			return nil, err
		}

		var result []model.Node
		for _, s := range successors {
			if _, found := seen[s]; !found {
				d.logger.Debugf("found link %s for collection %s", s, node.Address())
				childNode, err := remote.LoadCollection(ctx, s)
				if err != nil {
					return nil, err
				}
				result = append(result, &childNode)
				seen[s] = struct{}{}
			}
		}
		return result, nil
	})

	if err := tracker.Walk(ctx, handler, root); err != nil {
		return nil, err
	}

	return allDescs, nil
}

// getSuccessors retrieves all referenced collections from a source collection.
func getSuccessors(ctx context.Context, reference string, client registryclient.Remote) ([]string, error) {
	_, manBytes, err := client.GetManifest(ctx, reference)
	if err != nil {
		return nil, err
	}
	defer manBytes.Close()
	return ocimanifest.ResolveCollectionLinks(manBytes)
}
