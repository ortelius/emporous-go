package cache

import (
	"context"
	"io"
	"sync"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
)

// TODO(jpower432): Reimplement for UOR needs. Pulled from https://github.com/oras-project/oras/blob/main/internal/cache/target.go
// to use with oras.Copy command.

type closer func() error

func (fn closer) Close() error {
	return fn()
}

// Cache target struct.
type target struct {
	oras.Target
	cache content.Storage
}

// New generates a new target storage with caching.
func New(source oras.Target, cache content.Storage) oras.Target {
	return &target{
		Target: source,
		cache:  cache,
	}
}

// Fetch fetches the content identified by the descriptor.
func (p *target) Fetch(ctx context.Context, target ocispec.Descriptor) (io.ReadCloser, error) {
	rc, err := p.cache.Fetch(ctx, target)
	if err == nil {
		return rc, nil
	}

	rc, err = p.Target.Fetch(ctx, target)
	if err != nil {
		return nil, err
	}
	pr, pw := io.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)
	var pushErr error
	go func() {
		defer wg.Done()
		pushErr = p.cache.Push(ctx, target, pr)
	}()
	c := closer(func() error {
		rcErr := rc.Close()
		if err := pw.Close(); err != nil {
			return err
		}
		wg.Wait()
		if pushErr != nil {
			return pushErr
		}
		return rcErr
	})

	return struct {
		io.Reader
		io.Closer
	}{
		Reader: io.TeeReader(rc, pw),
		Closer: c,
	}, nil
}

// Exists returns true if the described content exists.
func (p *target) Exists(ctx context.Context, desc ocispec.Descriptor) (bool, error) {
	exists, err := p.cache.Exists(ctx, desc)
	if err == nil && exists {
		return true, nil
	}
	return p.Target.Exists(ctx, desc)
}
