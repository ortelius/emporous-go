package main

import (
	"github.com/spf13/cobra"

	"github.com/emporous/emporous-go/cmd/client/commands"
)

func main() {
	rootCmd := commands.NewRootCmd()
	cobra.CheckErr(rootCmd.Execute())
}
