package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

func NewCancelCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "Cancel job",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			job, _ := cmd.Flags().GetString("job")
			cli.XCancel(job, args...)
		},
	}

	var job string
	cmd.Flags().StringVarP(&job, "job", "j", "", "Cancel <JOB_ID>")
	return cmd
}
