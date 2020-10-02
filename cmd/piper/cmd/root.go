// Package cmd contains Cobra style commands.
package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCommand returns the root command.
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "piper",
		Short: "Piper short description",
		Long:  `Piper long description`,
	}

	root.AddCommand(newGenCommand())

	return root
}
