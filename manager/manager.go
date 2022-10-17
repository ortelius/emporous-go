package manager

import (
	"context"

	clientapi "github.com/uor-framework/uor-client-go/api/client/v1alpha1"
	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/registryclient"
	"github.com/uor-framework/uor-client-go/util/workspace"
)

// Manager defines methods for building, publishing, and retrieving UOR collections.
type Manager interface {
	// Build builds collection from input and store it in the underlying content store.
	// If successful, the root descriptor is returned.
	Build(ctx context.Context, source workspace.Workspace, config clientapi.DataSetConfiguration, destination string, client registryclient.Client) (string, error)
	// Push pushes collection to a remote location from the underlying content store.
	// If successful, the root descriptor is returned.
	Push(ctx context.Context, destination string, remote registryclient.Remote) (string, error)
	// Pull pulls a single collection to a specified storage destination.
	// If successful, the file locations are returned.
	Pull(ctx context.Context, source string, remote registryclient.Remote, destination content.Store) ([]string, error)
	// PullAll pulls linked collection to a specified storage destination.
	// If successful, the file locations are returned.
	// PullAll is similar to Pull with the exception that it walks a graph of linked collections
	// starting with the source collection reference.
	PullAll(ctx context.Context, source string, remote registryclient.Remote, destination content.Store) ([]string, error)
	// Update adds and removes content from a collection and stores the collection in the
	// underlying content store. If successful, the root descriptor is returned.
	Update(ctx context.Context, space workspace.Workspace, src string, dest string, add bool, remove bool, client registryclient.Client) (string, error)
}
