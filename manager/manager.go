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
	Build(ctx context.Context, source workspace.Workspace, config clientapi.DataSetConfiguration, destination string, client registryclient.Client) (string, error)
	// Push pushes collection to a remote location from the underlying content store.
	Push(ctx context.Context, destination string, remote registryclient.Remote) (string, error)
	// Pull pulls a single collection to a specified storage destination.
	Pull(ctx context.Context, source string, remote registryclient.Remote, destination content.Store) error
	// PullAll pulls linked collection to a specified storage destination.
	PullAll(ctx context.Context, source string, remote registryclient.Remote, destination content.Store) error
}
