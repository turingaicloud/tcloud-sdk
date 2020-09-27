package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

// TODO(), By default, no args
func NewConfigCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "User init workspace and download latest TACC config",
		Run: func(cmd *cobra.Command, args []string) {
			cli.XInit(args...)
		},
	}
}
