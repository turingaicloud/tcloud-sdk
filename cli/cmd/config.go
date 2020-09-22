package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

// TODO(), By default, no args
func NewConfigCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Configure user's account in tcloud CLI",
		Run: func(cmd *cobra.Command, args []string) {
			cli.XConfig(args...)
		},
	}
}
