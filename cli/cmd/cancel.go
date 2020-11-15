package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

func NewCancelCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "Cancel job",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			job, _ := cmd.Flags().GetString("job")
			cli.XCancel(job, args...)
		},
	}

	var job string
	cmd.Flags().StringVarP(&job, "job", "j", "", "Show <JOB_ID> status")
	return cmd
}
