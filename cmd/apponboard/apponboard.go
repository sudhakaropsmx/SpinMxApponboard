package apponboard

import (
	"github.com/spf13/cobra"
	"io"
	"github.com/sudhakaropsmx/spinmx/cmd/apponboard/pipeline"
)

type apponboardOptions struct{}

var (
	apponboardShort   = ""
	apponboardLong    = ""
	apponboardExample = ""
)

func NewAppOnboardCmd(out io.Writer) *cobra.Command {
	//options := apponboardOptions{}
	cmd := &cobra.Command{
		Use:     "apponboard",
		Aliases: []string{"apponboarding", "aob"},
		Short:   apponboardShort,
		Long:    apponboardLong,
		Example: apponboardExample,
	}  
	// create subcommands
	cmd.AddCommand(pipeline.NewPipelineCmd(out))
	return cmd
}
