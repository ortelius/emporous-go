package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// Version describes the version of the client
// set at build time or detected during runtime.
var Version string

// NewVersionCmd creates a new cobra.Command for the version subcommand.
func NewVersionCmd(rootOpts *RootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			if Version == "" {
				i, ok := debug.ReadBuildInfo()
				if ok {
					Version = i.Main.Version
				}
			}

			if Version == "" {
				_, _ = fmt.Fprintln(rootOpts.IOStreams.ErrOut, "could not determine build information")
			} else {
				_, _ = fmt.Fprintf(rootOpts.IOStreams.Out, "%s version: %v\n", filepath.Base(filepath.Base(os.Args[0])), Version)
			}
		},
	}
}
