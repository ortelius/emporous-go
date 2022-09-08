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

func (d DefaultManager) Pull(ctx context.Context, source string, remote registryclient.Remote, destination content.Store) error {
	_, err := d.pullCollection(ctx, source, destination, remote)
	if err != nil {
		return err
	}
	return nil
}

func (d DefaultManager) PullAll(ctx context.Context, source string, remote registryclient.Remote, destination content.Store) error {
	_, err := d.pullCollections(ctx, source, destination, remote)
	if err != nil {
		return err
	}
	return nil
}

// pullCollection pulls a single collection and returns the manifest descriptors and an error.
func (d DefaultManager) pullCollection(ctx context.Context, reference string, destination content.Store, remote registryclient.Remote) ([]ocispec.Descriptor, error) {
	desc, err := remote.Pull(ctx, reference, destination)
	if err != nil {
		if errors.Is(err, registryclient.ErrNoMatch) {
			d.logger.Infof("No matches found for collection %s", reference)
			return nil, nil
		}
		return nil, err
	}
	// Ensure the store is tagged with the new reference.
	return []ocispec.Descriptor{desc}, d.store.Tag(ctx, desc, reference)
}

// pullCollections pulls two or more collections and returns the manifest descriptors and an error.
func (d DefaultManager) pullCollections(ctx context.Context, source string, destination content.Store, remote registryclient.Remote) ([]ocispec.Descriptor, error) {
	root, err := remote.LoadCollection(ctx, source)
	if err != nil {
		return nil, err
	}
	return d.copy(ctx, &root, destination, remote)
}

// copy performs graph traversal of linked collections and performs collection copies filtered by the matcher.
func (d DefaultManager) copy(ctx context.Context, root model.Node, destination content.Store, remote registryclient.Remote) ([]ocispec.Descriptor, error) {
	seen := map[string]struct{}{}
	var manifestDesc []ocispec.Descriptor

	tracker := traversal.NewTracker(root, nil)
	handler := traversal.HandlerFunc(func(ctx context.Context, tracker traversal.Tracker, node model.Node) ([]model.Node, error) {

		descs, err := d.pullCollection(ctx, node.Address(), destination, remote)
		if err != nil {
			return nil, err
		}
		manifestDesc = append(manifestDesc, descs...)

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

	return manifestDesc, nil
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
