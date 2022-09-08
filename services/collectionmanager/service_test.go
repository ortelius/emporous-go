package collectionmanager

import (
	"context"
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
		pubAssertFunc func(string) bool
		workspace     string
		config        []byte
		resAssertFunc func(string) bool
		sev           managerapi.Diagnostic_Severity
		errMes        string
	}{
		{
			name:      "Success/ValidWorkspace",
			sev:       0,
			errMes:    "",
			workspace: "testdata/workspace",
			pubAssertFunc: func(s string) bool {
				return s == "sha256:2f0e884ddba718cba5eb540e3c0cb448ac0e72738a872be1618d839168b39032"
			},
			resAssertFunc: func(root string) bool {
				_, err := os.Stat(path.Join(root, "fish.jpg"))
				return err == nil
			},
		},
	}

	ctx := context.Background()

	testlogr, err := log.NewLogger(ioutil.Discard, "debug")
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
				Json:        c.config,
			}

			pResp, err := client.PublishContent(ctx, pRequest, opts...)
			if c.errMes != "" {
				require.EqualError(t, err, c.errMes)
			} else {
				require.NoError(t, err)
				require.True(t, c.pubAssertFunc(pResp.Digest))
			}

			destination := t.TempDir()
			rRequest := &managerapi.Retrieve_Request{
				Source:      fmt.Sprintf("%s/test:latest", u.Host),
				Destination: destination,
			}

			rResp, err := client.RetrieveContent(ctx, rRequest, opts...)
			if c.errMes != "" {
				rResp.Diagnostics[0].Severity = c.sev
				require.EqualError(t, err, c.errMes)
			} else {
				require.NoError(t, err)
				require.True(t, c.resAssertFunc(destination))
			}
		})
	}
}

var _ content.AttributeStore = testContentStore{}

type testContentStore struct {
	content.Store
}

func (t testContentStore) ResolveByAttribute(ctx context.Context, s string, matcher model.Matcher) ([]ocispec.Descriptor, error) {
	return nil, nil
}

func (t testContentStore) AttributeSchema(ctx context.Context, s string) (ocispec.Descriptor, error) {
	return ocispec.Descriptor{}, nil
}
