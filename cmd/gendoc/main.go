package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"github.com/uor-framework/uor-client-go/cmd/client/commands"
)

func main() {
	genDocCmd := &cobra.Command{
		Use:          "gendoc",
		Short:        "Generate UOR client CLI docs",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return doc.GenMarkdownTree(commands.NewRootCmd(), args[0])
		},
	}
	cobra.CheckErr(genDocCmd.Execute())
}
