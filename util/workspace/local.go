package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type localWorkspace struct {
	fs  afero.Fs
	dir string
}

var _ Workspace = &localWorkspace{}

// NewLocalWorkspace returns a new local workspace.
func NewLocalWorkspace(dir string) (Workspace, error) {
	// Get absolute path for provided dir
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	w := localWorkspace{
		dir: absDir,
	}
	return &w, w.init()
}

func (w *localWorkspace) init() error {
	if w.fs == nil {
		w.fs = afero.NewOsFs()
	}

	if err := w.fs.MkdirAll(w.dir, 0750); err != nil {
		return err
	}

	// Use a basepath FS to obviate joining paths later.
	// Do this after creating the dir using the underlying fs
	// so b.dir is not created under the base (itself).
	w.fs = afero.NewBasePathFs(w.fs, w.dir)

	return nil
}

// ReadObject reads the provided object from disk.
// In this implementation, key is a file path.
func (w *localWorkspace) ReadObject(_ context.Context, path string, obj interface{}) error {

	data, err := afero.ReadFile(w.fs, path)
	if err != nil {
		return err
	}

	switch v := obj.(type) {
	case []byte:
		if len(v) < len(data) {
			return io.ErrShortBuffer
		}
		copy(v, data)
	case io.Writer:
		_, err = v.Write(data)
	default:
		err = json.Unmarshal(data, obj)
	}
	return err
}

// WriteObject writes the provided object to disk.
// In this implementation, key is a file path.
func (w *localWorkspace) WriteObject(ctx context.Context, path string, obj interface{}) error {

	writer, err := w.GetWriter(ctx, path)
	if err != nil {
		return err
	}
	defer writer.(io.WriteCloser).Close()

	var data []byte
	switch v := obj.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	case io.Reader:
		data, err = io.ReadAll(v)
	default:
		data, err = json.Marshal(obj)
	}
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

// GetWriter returns an os.File as a writer.
// In this implementation, key is a file path.
func (w *localWorkspace) GetWriter(_ context.Context, path string) (io.Writer, error) {

	// Create a child dirs necessary.
	if err := w.fs.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return nil, fmt.Errorf("error creating object child path: %v", err)
	}

	writer, err := w.fs.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return nil, fmt.Errorf("error opening object file: %v", err)
	}

	return writer, nil
}

// Walk traverses the workspace directory.
func (w *localWorkspace) Walk(walkFunc filepath.WalkFunc) error {
	return afero.Walk(w.fs, ".", walkFunc)
}

// Path generates a path of a file with the workspace directory.
func (w *localWorkspace) Path(elem ...string) string {
	complete := []string{w.dir}
	return filepath.Join(append(complete, elem...)...)
}

// NewDirectory creates a new workspace under the current workspace.
func (w *localWorkspace) NewDirectory(path string) (Workspace, error) {
	if err := w.fs.MkdirAll(path, 0750); err != nil {
		return nil, err
	}
	space := &localWorkspace{
		fs:  afero.NewBasePathFs(w.fs, path),
		dir: w.Path(path),
	}
	return space, nil
}

// DeleteDirectory will delete a directory under a workspace.
func (w *localWorkspace) DeleteDirectory(path string) error {
	return w.fs.RemoveAll(path)
}
