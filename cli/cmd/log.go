package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

func NewLogCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "Check submitted jobs' log",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			job, _ := cmd.Flags().GetString("job")
			downloadPath, _ := cmd.Flags().GetString("filepath")
			cli.XLog(job, downloadPath, args...)
		},
	}

	var job string
	cmd.Flags().StringVarP(&job, "job", "j", "", "Show <JOB_ID> log")
	var downloadPath string
	cmd.Flags().StringVarP(&downloadPath, "filepath", "f", "", "Download file in user directory")
	return cmd
}
