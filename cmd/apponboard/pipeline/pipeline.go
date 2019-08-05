package pipeline

import (
	"github.com/spf13/cobra"
	"io"
)

type pipelineOptions struct{}

var (
	pipelineShort   = ""
	pipelineLong    = ""
	pipelineExample = ""
)

func NewPipelineCmd(out io.Writer) *cobra.Command {
	options := pipelineOptions{}
	cmd := &cobra.Command{
		Use:     "pipeline",
		Aliases: []string{"pipelines", "pi"},
		Short:   pipelineShort,
		Long:    pipelineLong,
		Example: pipelineExample,
	}

	// create subcommands
	cmd.AddCommand(NewCreatePipelineCmd(options))
	return cmd
}
