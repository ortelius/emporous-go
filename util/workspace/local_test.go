package workspace

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestLocalWorkspace(t *testing.T) {

	underlyingFS := afero.NewMemMapFs()
	workspace := localWorkspace{
		fs:  underlyingFS,
		dir: filepath.Join("foo"),
	}
	require.NoError(t, workspace.init())

	ctx := context.Background()

	type object struct {
		SomeData string
	}

	inObj := object{
		SomeData: "bar",
	}
	testpath := "foo.txt"

	require.NoError(t, workspace.WriteObject(ctx, testpath, inObj))
	info, err := underlyingFS.Stat("foo/foo.txt")
	require.NoError(t, err)
	require.True(t, info.Mode().IsRegular())
	info, err = workspace.fs.Stat("foo.txt")
	require.NoError(t, err)
	require.True(t, info.Mode().IsRegular())

	var outObj object
	require.NoError(t, workspace.ReadObject(ctx, testpath, &outObj))
	require.Equal(t, inObj, outObj)

	subWorkspace, err := workspace.NewDirectory("anotherworkspace")
	require.NoError(t, err)

	require.NoError(t, subWorkspace.WriteObject(ctx, testpath, inObj))
	info, err = underlyingFS.Stat("foo/anotherworkspace/foo.txt")
	require.NoError(t, err)
	require.True(t, info.Mode().IsRegular())

	var outObj2 object
	require.NoError(t, subWorkspace.ReadObject(ctx, testpath, &outObj2))
	require.Equal(t, inObj, outObj)

	require.NoError(t, workspace.DeleteDirectory("anotherworkspace"))
	_, err = underlyingFS.Stat("foo/anotherworkspace")
	require.ErrorIs(t, err, os.ErrNotExist)
}
