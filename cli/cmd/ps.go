package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

// TODO(), By default, no args
func NewPSCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	return &cobra.Command{
		Use:   "ps",
		Short: "Check slurm jobs' status",
		Run: func(cmd *cobra.Command, args []string) {
			cli.XPS(args...)
		},
	}
}
