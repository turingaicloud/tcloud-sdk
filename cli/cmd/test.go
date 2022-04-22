package cmd

import (
	"tcloud-sdk/cli/tcloudcli"

	"github.com/spf13/cobra"
)

func NewTestCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "For test only",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cli.XTest(args...)
		},
	}
}
