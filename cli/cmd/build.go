package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

func NewBuildCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Parse tuxiv.confg and Setup conda environment",
		Run: func(cmd *cobra.Command, args []string) {
			cli.XBuild(args...)
		},
	}
}
