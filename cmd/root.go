package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/sudhakaropsmx/spinmx/cmd/apponboard"
	"github.com/sudhakaropsmx/spinmx/cmd/project"
	"github.com/spinnaker/spin/cmd/application"
	"github.com/spinnaker/spin/cmd/pipeline"
	"github.com/spinnaker/spin/cmd/pipeline-template"
	"github.com/sudhakaropsmx/spinmx/version"
	
)

type RootOptions struct {
	configFile       string
	GateEndpoint     string
	ignoreCertErrors bool
	quiet            bool
	color            bool
	outputFormat     string
	defaultHeaders   string
}

func Execute(out io.Writer) error {
	cmd := NewCmdRoot(out)
	return cmd.Execute()
}

func NewCmdRoot(out io.Writer) *cobra.Command {
	options := RootOptions{}

	cmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version.String(),
	}

	cmd.PersistentFlags().StringVar(&options.configFile, "config", "", "path to config file (default $HOME/.spin/config)")
	cmd.PersistentFlags().StringVar(&options.GateEndpoint, "gate-endpoint", "", "Gate (API server) endpoint (default http://localhost:8084)")
	cmd.PersistentFlags().BoolVarP(&options.ignoreCertErrors, "insecure", "k", false, "ignore certificate errors")
	cmd.PersistentFlags().BoolVarP(&options.quiet, "quiet", "q", false, "squelch non-essential output")
	cmd.PersistentFlags().BoolVar(&options.color, "no-color", true, "disable color")
	cmd.PersistentFlags().StringVar(&options.outputFormat, "output", "", "configure output formatting")
	cmd.PersistentFlags().StringVar(&options.defaultHeaders, "default-headers", "", "configure default headers for gate client as comma separated list (e.g. key1=value1,key2=value2)")

	// create subcommands
	cmd.AddCommand(apponboard.NewAppOnboardCmd(out))
	cmd.AddCommand(application.NewApplicationCmd(out))
	cmd.AddCommand(pipeline.NewPipelineCmd(out))
	cmd.AddCommand(pipeline_template.NewPipelineTemplateCmd(out))
	cmd.AddCommand(project.NewProjectCmd(out))
	
	return cmd
}
