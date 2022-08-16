package main

import (
	"github.com/spf13/cobra"

	"github.com/uor-framework/uor-client-go/cli"
)

func main() {
	rootCmd := cli.NewRootCmd()
	cobra.CheckErr(rootCmd.Execute())
}
