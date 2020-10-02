package cmd

import (
	"github.com/finderseyes/piper/pipes"
	"github.com/finderseyes/piper/pipes/io"
	"github.com/spf13/cobra"
)

func newGenCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "gen",
		Short: "Generate pipes",
		Long:  `Generate pipes`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			factory := io.NewFileWriterFactory()
			// Do Stuff Here
			generator := pipes.NewGenerator(args[0],
				pipes.WithWriterFactory(factory),
			)
			return generator.Execute()
		},
	}
}
