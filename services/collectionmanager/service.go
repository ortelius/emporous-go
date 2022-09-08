package collectionmanager

import (
	"context"

	"oras.land/oras-go/v2/content/file"

	"github.com/uor-framework/uor-client-go/api/client/v1alpha1"
	managerapi "github.com/uor-framework/uor-client-go/api/services/collectionmanager/v1alpha1"
	"github.com/uor-framework/uor-client-go/attributes/matchers"
	"github.com/uor-framework/uor-client-go/config"
	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/manager"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient"
	"github.com/uor-framework/uor-client-go/util/workspace"
)

var _ managerapi.CollectionManagerServer = &service{}

type service struct {
	managerapi.UnimplementedCollectionManagerServer
	mg      manager.Manager
	options ServiceOptions
}

// ServiceOptions configure the collection router service with default remote
// and collection caching options.
type ServiceOptions struct {
	Insecure    bool
	PlainHTTP   bool
	AuthConfigs []string
	PullCache   content.Store
}

// FromManager returns a CollectionManager API server from a Manager type.
func FromManager(mg manager.Manager, serviceOptions ServiceOptions) managerapi.CollectionManagerServer {
	return &service{
		mg:      mg,
		options: serviceOptions,
	}
}

// PublishContent publishes collection content to a storage provide based on client input.
func (s *service) PublishContent(ctx context.Context, message *managerapi.Publish_Request) (*managerapi.Publish_Response, error) {
	client, err := orasclient.NewClient(
		orasclient.WithCache(s.options.PullCache),
		orasclient.WithPlainHTTP(s.options.PlainHTTP),
		orasclient.WithAuthConfigs(s.options.AuthConfigs),
		orasclient.SkipTLSVerify(s.options.Insecure))
	if err != nil {
		return &managerapi.Publish_Response{
			Diagnostics: []*managerapi.Diagnostic{
				{
					Severity: 1,
					Summary:  "PublishError",
					Detail:   err.Error(),
				},
			},
		}, err
	}

	space, err := workspace.NewLocalWorkspace(message.Source)
	if err != nil {
		return &managerapi.Publish_Response{
			Diagnostics: []*managerapi.Diagnostic{
				{
					Severity: 1,
					Summary:  "PublishError",
					Detail:   err.Error(),
				},
			},
		}, err
	}

	var dsConfig v1alpha1.DataSetConfiguration
	if message.Json != nil {
		dsConfig, err = config.LoadDataSetConfig(message.Json)
		if err != nil {
			return &managerapi.Publish_Response{
				Diagnostics: []*managerapi.Diagnostic{
					{
						Severity: 1,
						Summary:  "PublishError",
						Detail:   err.Error(),
					},
				},
			}, err
		}
	}

	_, err = s.mg.Build(ctx, space, dsConfig, message.Destination, client)
	if err != nil {
		return &managerapi.Publish_Response{
			Diagnostics: []*managerapi.Diagnostic{
				{
					Severity: 1,
					Summary:  "PublishError",
					Detail:   err.Error(),
				},
			},
		}, err
	}

	digest, err := s.mg.Push(ctx, message.Destination, client)
	if err != nil {
		return &managerapi.Publish_Response{
			Diagnostics: []*managerapi.Diagnostic{
				{
					Severity: 1,
					Summary:  "PublishError",
					Detail:   err.Error(),
				},
			},
		}, err
	}

	return &managerapi.Publish_Response{Digest: digest}, nil
}

// RetrieveContent retrieves collection contact from a storage provider based on client input.
func (s *service) RetrieveContent(ctx context.Context, message *managerapi.Retrieve_Request) (*managerapi.Retrieve_Response, error) {
	attrSet, err := config.ConvertToModel(message.Filter.AsMap())
	if err != nil {
		return &managerapi.Retrieve_Response{
			Diagnostics: []*managerapi.Diagnostic{
				{
					Severity: 1,
					Summary:  "RetrieveError",
					Detail:   err.Error(),
				},
			},
		}, err
	}

	var matcher matchers.PartialAttributeMatcher = attrSet.List()
	client, err := orasclient.NewClient(
		orasclient.WithCache(s.options.PullCache),
		orasclient.WithAuthConfigs(s.options.AuthConfigs),
		orasclient.WithPlainHTTP(s.options.PlainHTTP),
		orasclient.SkipTLSVerify(s.options.Insecure),
		orasclient.WithPullableAttributes(matcher),
	)
	if err != nil {
		return &managerapi.Retrieve_Response{
			Diagnostics: []*managerapi.Diagnostic{
				{
					Severity: 1,
					Summary:  "RetrieveError",
					Detail:   err.Error(),
				},
			},
		}, err
	}
	err = s.mg.PullAll(ctx, message.Source, client, file.New(message.Destination))
	if err != nil {
		return &managerapi.Retrieve_Response{
			Diagnostics: []*managerapi.Diagnostic{
				{
					Severity: 1,
					Summary:  "RetrieveError",
					Detail:   err.Error(),
				},
			},
		}, err
	}

	return &managerapi.Retrieve_Response{}, nil
}
