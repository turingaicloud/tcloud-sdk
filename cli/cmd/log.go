package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

func NewLogCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "Check submitted jobs' log",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			job, _ := cmd.Flags().GetString("job")
			cli.XLog(job, args...)
		},
	}

	var job string
	cmd.Flags().StringVarP(&job, "job", "j", "", "Show <JOB_ID> status")
	return cmd
}
