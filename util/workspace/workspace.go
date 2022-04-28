package workspace

import (
	"context"
	"io"
)

// Workspace defines methods for accessing and publishing
// files in a local or remote context.
type Workspace interface {
	ReadObject(context.Context, string, interface{}) error
	WriteObject(context.Context, string, interface{}) error
	GetWriter(context.Context, string) (io.Writer, error)
	Open(context.Context, string) (io.ReadCloser, error)
}
