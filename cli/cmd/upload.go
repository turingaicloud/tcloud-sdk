package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

func NewUploadCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload file to TACC_USERDIR",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cli.XUpload(args...)
		},
	}

	return cmd
}
