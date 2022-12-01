package layout

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"
	orascontent "oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/errdef"

	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/model/traversal"
	"github.com/uor-framework/uor-client-go/nodes/collection"
	"github.com/uor-framework/uor-client-go/nodes/collection/loader"
	"github.com/uor-framework/uor-client-go/nodes/descriptor/v2"
)

var (
	_ content.Store          = &Layout{}
	_ content.GraphStore     = &Layout{}
	_ content.AttributeStore = &Layout{}
)

const indexFile = "index.json"

// Layout implements the storage interface by wrapping the oras
// content.Storage.
type Layout struct {
	internal orascontent.Storage
	resolver sync.Map // map[string]ocispec.Descriptor
	graph    *collection.Collection
	index    *ocispec.Index
	rootPath string
	mu       sync.Mutex
}

// New initializes a new local file store in an OCI layout format.
func New(rootPath string) (*Layout, error) {
	return NewWithContext(context.Background(), rootPath)
}

// NewWithContext initializes a new local file store in an OCI layout format.
func NewWithContext(ctx context.Context, rootPath string) (*Layout, error) {
	l := &Layout{
		internal: oci.NewStorage(rootPath),
		resolver: sync.Map{},
		graph:    collection.New(rootPath),
		rootPath: filepath.Clean(rootPath),
	}

	return l, l.init(ctx)
}

// init performs initial layout checks and loads the index.
func (l *Layout) init(ctx context.Context) error {
	if err := l.validateOCILayoutFile(); err != nil {
		return err
	}
	return l.loadIndex(ctx)
}

// Fetch fetches the content identified by the descriptor.
func (l *Layout) Fetch(ctx context.Context, desc ocispec.Descriptor) (io.ReadCloser, error) {
	return l.internal.Fetch(ctx, desc)
}

// Push pushes the content, matching the expected descriptor.
func (l *Layout) Push(ctx context.Context, desc ocispec.Descriptor, content io.Reader) error {
	if err := l.internal.Push(ctx, desc, content); err != nil {
		return err
	}

	fetcherFn := func(ctx context.Context, desc ocispec.Descriptor) ([]byte, error) {
		return orascontent.FetchAll(ctx, l, desc)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return loader.AddManifest(ctx, l.graph, fetcherFn, desc)
}

// Exists returns whether a descriptor exits in the file store.
func (l *Layout) Exists(ctx context.Context, desc ocispec.Descriptor) (bool, error) {
	return l.internal.Exists(ctx, desc)
}

// Resolve resolves a reference to a descriptor.
func (l *Layout) Resolve(_ context.Context, reference string) (ocispec.Descriptor, error) {
	desc, ok := l.resolver.Load(reference)
	if !ok {
		return ocispec.Descriptor{}, &content.ErrNotStored{Reference: reference}
	}
	return desc.(ocispec.Descriptor), nil
}

// Predecessors returns the nodes directly pointing to the current node.
func (l *Layout) Predecessors(_ context.Context, node ocispec.Descriptor) ([]ocispec.Descriptor, error) {
	var predecessors []ocispec.Descriptor
	nodes := l.graph.To(node.Digest.String())
	for _, n := range nodes {
		desc, ok := n.(*v2.Node)
		if ok {
			predecessors = append(predecessors, desc.Descriptor())
		}
	}
	return predecessors, nil
}

// ResolveByAttribute returns descriptors linked to the reference that satisfy the specified matcher.
// Matcher is expected to compare attributes of nodes to set criteria. If the matcher is nil, return values
// are nil.
func (l *Layout) ResolveByAttribute(ctx context.Context, reference string, matcher model.Matcher) ([]ocispec.Descriptor, error) {
	if matcher == nil {
		return nil, nil
	}

	var res []ocispec.Descriptor
	desc, err := l.Resolve(ctx, reference)
	if err != nil {
		return nil, err
	}

	root := l.graph.NodeByID(desc.Digest.String())
	if root == nil {
		return nil, fmt.Errorf("node %q does not exist in graph", reference)
	}

	tracker := traversal.NewTracker(root, nil)
	handler := traversal.HandlerFunc(func(ctx context.Context, tracker traversal.Tracker, node model.Node) ([]model.Node, error) {
		match, err := matcher.Matches(node)
		if err != nil {
			return nil, err
		}
		if match {
			desc, ok := node.(*v2.Node)
			if ok {

				// Check that the blob actually exists within the file
				// store. This will filter out blobs in the event that this is a
				// sparse manifest.
				exists, err := l.internal.Exists(ctx, desc.Descriptor())
				if err != nil {
					return nil, err
				}
				if exists {
					res = append(res, desc.Descriptor())
				}
			}
		}

		return l.graph.From(node.ID()), nil
	})

	if err := tracker.Walk(ctx, handler, root); err != nil {
		return nil, err
	}

	return res, err
}

// AttributeSchema returns the descriptor containing the given attribute schema for a given reference.
func (l *Layout) AttributeSchema(ctx context.Context, reference string) (ocispec.Descriptor, error) {
	desc, err := l.Resolve(ctx, reference)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	root := l.graph.NodeByID(desc.Digest.String())
	if root == nil {
		return ocispec.Descriptor{}, fmt.Errorf("node %q does not exist in graph", reference)
	}
	var res ocispec.Descriptor
	var stopErr = errors.New("stop")
	tracker := traversal.NewTracker(root, nil)
	handler := traversal.HandlerFunc(func(ctx context.Context, tracker traversal.Tracker, node model.Node) ([]model.Node, error) {
		desc, ok := node.(*v2.Node)
		if ok {
			if desc.Descriptor().MediaType == uorspec.MediaTypeSchemaDescriptor {
				res = desc.Descriptor()
				return nil, stopErr
			}
		}
		return l.graph.From(node.ID()), nil
	})

	err = tracker.Walk(ctx, handler, root)
	if err == nil {
		return ocispec.Descriptor{}, fmt.Errorf("reference %s is not a schema address", reference)
	}

	if err != nil && !errors.Is(err, stopErr) {
		return ocispec.Descriptor{}, err
	}

	return res, nil
}

// Tag tags a descriptor with a reference string.
// A reference should be either a valid tag (e.g. "latest"),
// or a digest matching the descriptor (e.g. "@sha256:abc123").
func (l *Layout) Tag(ctx context.Context, desc ocispec.Descriptor, reference string) error {
	if err := validateReference(reference); err != nil {
		return err
	}

	exists, err := l.Exists(ctx, desc)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%s: %s: %w", desc.Digest, desc.MediaType, errdef.ErrNotFound)
	}

	if desc.Annotations == nil {
		desc.Annotations = map[string]string{}
	}
	desc.Annotations[ocispec.AnnotationRefName] = reference

	l.resolver.Store(reference, desc)

	return l.SaveIndex()
}

// Index returns an index manifest object.
func (l *Layout) Index() (ocispec.Index, error) {
	return *l.index, nil
}

// SaveIndex writes the index.json to the file system
func (l *Layout) SaveIndex() error {
	// first need to update the index
	var descs []ocispec.Descriptor
	l.resolver.Range(func(key, value interface{}) bool {
		desc := value.(ocispec.Descriptor)
		if desc.Annotations == nil {
			desc.Annotations = map[string]string{}
		}
		desc.Annotations[ocispec.AnnotationRefName] = key.(string)
		descs = append(descs, desc)
		return true
	})
	l.index.Manifests = descs
	indexJSON, err := json.Marshal(l.index)
	if err != nil {
		return err
	}
	path := filepath.Join(l.rootPath, indexFile)
	return ioutil.WriteFile(path, indexJSON, 0640)
}

// loadIndex loads all information from the index.json
// into the resolver and graph.
func (l *Layout) loadIndex(ctx context.Context) error {
	path := filepath.Join(l.rootPath, indexFile)
	indexFile, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		l.index = &ocispec.Index{
			Versioned: specs.Versioned{
				SchemaVersion: 2,
			},
		}

		return nil
	}
	defer indexFile.Close()

	fetcherFn := func(ctx context.Context, desc ocispec.Descriptor) ([]byte, error) {
		return orascontent.FetchAll(ctx, l, desc)
	}

	if err := json.NewDecoder(indexFile).Decode(&l.index); err != nil {
		return err
	}

	for _, d := range l.index.Manifests {
		key, ok := d.Annotations[ocispec.AnnotationRefName]
		if ok {
			l.resolver.Store(key, d)
		}

		if err := l.loadReference(ctx, fetcherFn, d); err != nil {
			return err
		}
	}

	return nil
}

func (l *Layout) loadReference(ctx context.Context, fetcherFn loader.FetcherFunc, manifest ocispec.Descriptor) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return loader.LoadFromManifest(ctx, l.graph, fetcherFn, manifest)
}

// validateOCILayoutFile ensure the 'oci-layout' file exists in the
// root directory and contains a valid version.
func (l *Layout) validateOCILayoutFile() error {
	layoutFilePath := filepath.Join(l.rootPath, ocispec.ImageLayoutFile)
	layoutFile, err := os.Open(layoutFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to open OCI layout file: %w", err)
		}

		layout := ocispec.ImageLayout{
			Version: ocispec.ImageLayoutVersion,
		}
		layoutJSON, err := json.Marshal(layout)
		if err != nil {
			return fmt.Errorf("failed to marshal OCI layout file: %w", err)
		}

		return ioutil.WriteFile(layoutFilePath, layoutJSON, 0666)
	}
	defer layoutFile.Close()

	var layout *ocispec.ImageLayout
	err = json.NewDecoder(layoutFile).Decode(&layout)
	if err != nil {
		return fmt.Errorf("failed to decode OCI layout file: %w", err)
	}
	if layout.Version != ocispec.ImageLayoutVersion {
		return errdef.ErrUnsupportedVersion
	}

	return nil
}

// validateReference ensures the build reference
// contains a tag component.
func validateReference(name string) error {
	parts := strings.SplitN(name, "/", 2)
	if len(parts) == 1 {
		return fmt.Errorf("reference %q: missing repository", name)
	}
	path := parts[1]
	if index := strings.Index(path, "@"); index != -1 {
		return fmt.Errorf("%q: %w", name, errdef.ErrInvalidReference)
	} else if index := strings.Index(path, ":"); index != -1 {
		// tag found
		return nil
	} else {
		// empty reference
		return fmt.Errorf("reference %q: missing tag component", name)
	}
}
