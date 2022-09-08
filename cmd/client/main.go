package main

import (
	"github.com/spf13/cobra"

	"github.com/uor-framework/uor-client-go/cmd/client/commands"
)

func main() {
	rootCmd := commands.NewRootCmd()
	cobra.CheckErr(rootCmd.Execute())
}
