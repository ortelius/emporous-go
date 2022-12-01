package defaultmanager

import (
	"context"

	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/registryclient"
)

// Pull pulls a single collection to a specified storage destination.
// If successful, the file locations are returned.
func (d DefaultManager) Pull(ctx context.Context, source string, remote registryclient.Remote, destination content.Store) ([]string, error) {
	rootDesc, descs, err := remote.Pull(ctx, source, destination)
	if err != nil {
		return nil, err
	}

	// Ensure the store is tagged with the new reference.
	if len(rootDesc.Digest) != 0 {
		if err := d.store.Tag(ctx, rootDesc, source); err != nil {
			return nil, err
		}
	}

	var digests []string
	for _, desc := range descs {
		digests = append(digests, desc.Digest.String())
		d.logger.Infof("Found matching digest %s", desc.Digest)
	}
	return digests, nil
}

// PullAll pulls linked collection to a specified storage destination.
// If successful, the file locations are returned.
// PullAll is similar to Pull with the exception that it walks a graph of linked collections
// starting with the source collection reference.
func (d DefaultManager) PullAll(ctx context.Context, source string, remote registryclient.Remote, destination content.Store) ([]string, error) {
	descs, err := remote.PullWithLinks(ctx, source, destination)
	if err != nil {
		return nil, err
	}

	var digests []string
	for _, desc := range descs {
		digests = append(digests, desc.Digest.String())
		d.logger.Infof("Found matching digest %s", desc.Digest)
	}
	return digests, nil
}
