package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

// TODO(), By default, no args
func NewSubmitCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	return &cobra.Command{
		Use:   "submit",
		Short: "Submit a slurm job",
		Run: func(cmd *cobra.Command, args []string) {
			cli.XSubmit(args...)
		},
	}
}
