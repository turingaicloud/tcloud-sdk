package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

// TODO(), By default, no args
func NewAttachCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	return &cobra.Command{
		Use:   "attach",
		Short: "Attach user's slurm command in tcloud CLI",
		Run: func(cmd *cobra.Command, args []string) {
			cli.XAttach(args...)
		},
	}
}
