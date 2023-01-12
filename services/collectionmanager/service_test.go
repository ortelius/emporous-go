package collectionmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/google/go-containerregistry/pkg/registry"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/structpb"
	"oras.land/oras-go/v2/content/memory"

	managerapi "github.com/uor-framework/uor-client-go/api/services/collectionmanager/v1alpha1"
	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/log"
	"github.com/uor-framework/uor-client-go/manager/defaultmanager"
	"github.com/uor-framework/uor-client-go/model"
)

func dialer(srv managerapi.CollectionManagerServer) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	managerapi.RegisterCollectionManagerServer(server, srv)

	go func() {
		if err := server.Serve(listener); err != nil {
			fmt.Println(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestCollectionManagerServer_All(t *testing.T) {
	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	cases := []struct {
		name          string
		pubAssertFunc func(*managerapi.Publish_Response) bool
		workspace     string
		collection    map[string]map[string]interface{}
		filter        json.RawMessage
		resAssertFunc func(*managerapi.Retrieve_Response, string) bool
		sev           managerapi.Diagnostic_Severity
		errMes        string
	}{
		{
			name:      "Success/ValidWorkspace",
			workspace: "testdata/workspace",
			resAssertFunc: func(_ *managerapi.Retrieve_Response, root string) bool {
				_, err := os.Stat(path.Join(root, "fish.jpg"))
				return err == nil
			},
		},
		{
			name:      "Success/WithConfig",
			workspace: "testdata/workspace",
			collection: map[string]map[string]interface{}{
				"*.jpg": {
					"animal": true,
				},
			},
			filter: []byte(`{"unknown":{"animal":true}}`),
			resAssertFunc: func(_ *managerapi.Retrieve_Response, root string) bool {
				_, err := os.Stat(path.Join(root, "fish.jpg"))
				return err == nil
			},
		},
		{
			name:      "Warning/FilteredCollection",
			sev:       2,
			errMes:    "",
			filter:    []byte(`{"test": "test"}`),
			workspace: "testdata/workspace",
			resAssertFunc: func(resp *managerapi.Retrieve_Response, _ string) bool {
				return len(resp.Diagnostics) != 0 && resp.Diagnostics[0].Severity == 2
			},
		},
	}

	ctx := context.Background()

	testlogr, err := log.NewLogrusLogger(ioutil.Discard, "debug")
	require.NoError(t, err)

	manager := defaultmanager.New(testContentStore{Store: memory.New()}, testlogr)
	srv := FromManager(manager, ServiceOptions{PlainHTTP: true})

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(srv)))
	require.NoError(t, err)
	defer conn.Close()

	client := managerapi.NewCollectionManagerClient(conn)
	var opts []grpc.CallOption

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			pRequest := &managerapi.Publish_Request{
				Source:      c.workspace,
				Destination: fmt.Sprintf("%s/test:latest", u.Host),
			}

			if c.collection != nil {
				var files []*managerapi.File
				for file, attr := range c.collection {
					a, err := structpb.NewStruct(attr)
					require.NoError(t, err)
					f := &managerapi.File{
						File:       file,
						Attributes: a,
					}
					files = append(files, f)
				}

				pRequest.Collection = &managerapi.Collection{
					Files: files,
				}
			}

			pResp, err := client.PublishContent(ctx, pRequest, opts...)
			if c.errMes != "" {
				require.EqualError(t, err, c.errMes)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, err)
			destination := t.TempDir()
			rRequest := &managerapi.Retrieve_Request{
				Source:      fmt.Sprintf("%s/test@%s", u.Host, pResp.Digest),
				Destination: destination,
			}

			if c.filter != nil {
				var data map[string]interface{}
				err = json.Unmarshal(c.filter, &data)
				require.NoError(t, err)
				filter, err := structpb.NewStruct(data)
				require.NoError(t, err)
				rRequest.Filter = filter
			}

			rResp, err := client.RetrieveContent(ctx, rRequest, opts...)
			if c.errMes != "" {
				require.EqualError(t, err, c.errMes)
			} else {
				require.NoError(t, err)
				require.True(t, c.resAssertFunc(rResp, destination))
			}
		})
	}
}

var _ content.AttributeStore = testContentStore{}

type testContentStore struct {
	content.Store
}

func (t testContentStore) ResolveByAttribute(_ context.Context, _ string, _ model.Matcher) ([]ocispec.Descriptor, error) {
	return nil, nil
}

func (t testContentStore) AttributeSchema(_ context.Context, _ string) (ocispec.Descriptor, error) {
	return ocispec.Descriptor{}, nil
}
