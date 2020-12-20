package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

func NewDownloadCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download file from TACC",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			IsDir, _ := cmd.Flags().GetBool("recursive")
			cli.XDownload(IsDir, args...)
		},
	}

	var IsDir bool
	cmd.Flags().BoolVarP(&IsDir, "recursive", "r", false, "recursively download")
	return cmd
}
