// Package cmd contains Cobra style commands.
package cmd

import (
	"github.com/finderseyes/piper/pipes"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "piper",
		Short: "Piper short description",
		Long: `Piper long description`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Do Stuff Here
			generator := pipes.NewGenerator(args[0])
			return generator.Execute()
		},
	}
}
