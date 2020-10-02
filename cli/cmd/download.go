package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

func NewDownloadCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	return &cobra.Command{
		Use:   "download",
		Short: "Download repo from <URL>",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cli.XDownload(args...)
		},
	}
}
