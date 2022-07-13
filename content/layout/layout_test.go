package layout

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"
)

func TestExists(t *testing.T) {
	cacheDir := "testdata/valid"
	l, err := New(cacheDir)
	require.NoError(t, err)
	type spec struct {
		name     string
		desc     ocispec.Descriptor
		expRes   bool
		expError string
	}

	cases := []spec{
		{
			name: "Success/DescExists",
			desc: ocispec.Descriptor{
				MediaType: "application/vnd.uor.config.v1+json",
				Size:      2,
				Digest:    "sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
			},
			expRes: true,
		},
		{
			name: "Success/DescDoesNotExist",
			desc: ocispec.Descriptor{
				MediaType: "application/vnd.uor.config.v1+json",
				Size:      3,
				Digest:    "sha256:44136fa356b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
			},
			expRes: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := l.Exists(context.TODO(), c.desc)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expRes, res)
			}
		})
	}
}

func TestTag(t *testing.T) {
	cacheDir := t.TempDir()
	l, err := New(cacheDir)
	require.NoError(t, err)

	l.descriptorLookup.Store("test", ocispec.Descriptor{})
	require.NoError(t, l.SaveIndex())

	_, err = os.Stat(filepath.Join(cacheDir, indexFile))
	require.NoError(t, err)

	ii, err := l.Index()
	require.NoError(t, err)
	require.Len(t, ii.Manifests, 1)
	desc := ii.Manifests[0]
	refName := desc.Annotations[ocispec.AnnotationRefName]

	require.Equal(t, "test", refName)
}

func TestSaveIndex(t *testing.T) {
	cacheDir := t.TempDir()
	l, err := New(cacheDir)
	require.NoError(t, err)

	l.descriptorLookup.Store("test", ocispec.Descriptor{})
	require.NoError(t, l.SaveIndex())

	_, err = os.Stat(filepath.Join(cacheDir, indexFile))
	require.NoError(t, err)

	ii, err := l.Index()
	require.NoError(t, err)
	require.Len(t, ii.Manifests, 1)
	desc := ii.Manifests[0]
	refName := desc.Annotations[ocispec.AnnotationRefName]

	require.Equal(t, "test", refName)
}

func TestResolve(t *testing.T) {
	cacheDir := "testdata/valid"
	l, err := New(cacheDir)
	require.NoError(t, err)
	type spec struct {
		name     string
		ref      string
		expDesc  ocispec.Descriptor
		expError string
	}

	cases := []spec{
		{
			name: "Success/RefExists",
			ref:  "localhost:5001/test:latest",
			expDesc: ocispec.Descriptor{
				MediaType:   "application/vnd.oci.image.manifest.v1+json",
				Size:        1013,
				Digest:      "sha256:473f7d69dbc51105aff4bb2f7ec80e27402d2f40c3e9a076e8c773b15969eadf",
				Annotations: map[string]string{"org.opencontainers.image.ref.name": "localhost:5001/test:latest"},
			},
		},
		{
			name:     "Failure/RefDoesNotExist",
			ref:      "localhost:5001/notexists:latest",
			expError: "descriptor for reference localhost:5001/notexists:latest is not stored",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			desc, err := l.Resolve(context.TODO(), c.ref)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expDesc, desc)
			}
		})
	}
}

func TestLoadIndex(t *testing.T) {
	cacheDir := "testdata/valid"
	l, err := New(cacheDir)
	require.NoError(t, err)

	require.NoError(t, l.loadIndex())

	ii, err := l.Index()
	require.NoError(t, err)

	require.Len(t, ii.Manifests, 1)
	require.Equal(t, "sha256:473f7d69dbc51105aff4bb2f7ec80e27402d2f40c3e9a076e8c773b15969eadf", ii.Manifests[0].Digest.String())
}

func TestValidateOCILayoutfile(t *testing.T) {
	type spec struct {
		name     string
		cacheDir string
		expError string
	}

	cases := []spec{
		{
			name:     "Success/OCILayoutFileWithCorrectVersion",
			cacheDir: "testdata/valid",
		},
		{
			name:     "Failure/WrongVersion",
			cacheDir: "testdata/invalid",
			expError: "unsupported version",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := &Layout{rootPath: c.cacheDir}
			err := l.validateOCILayoutFile()
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
