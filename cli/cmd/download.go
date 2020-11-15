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
			// job, _ := cmd.Flags().GetString("job")
			// downloadPath, _ := cmd.Flags().GetString("filepath")
			// cli.XLog(job, downloadPath, IsDir, args...)
			cli.XDownload(IsDir, args...)
		},
	}

	// var job string
	// cmd.Flags().StringVarP(&job, "job", "j", "", "Show <JOB_ID> log")
	var IsDir bool
	cmd.Flags().BoolVarP(&IsDir, "recursive", "r", false, "recursively download")
	// var downloadPath string
	// cmd.Flags().StringVarP(&downloadPath, "filepath", "f", "", "Download file in user directory")
	return cmd
}
