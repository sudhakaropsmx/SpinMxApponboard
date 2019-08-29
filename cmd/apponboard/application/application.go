package application

import (
	"github.com/spf13/cobra"
	"io"
)

type applicationOptions struct{}

var (
	applicationShort   = ""
	applicationLong    = ""
	applicationExample = ""
)

func NewApplicationCmd(out io.Writer) *cobra.Command {
	options := applicationOptions{}
	cmd := &cobra.Command{
		Use:     "application",
		Aliases: []string{"applications", "app"},
		Short:   applicationShort,
		Long:    applicationLong,
		Example: applicationExample,
	}

	// create subcommands
	cmd.AddCommand(NewCreateApplicationCmd(options))
	return cmd
}
