package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

// TODO(), By default, no args
func NewPSCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ps",
		Short: "Check slurm jobs' status",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			job, err := cmd.Flags().GetString("job")
			if err != nil {
				cli.XPS("", args...)
			} else {
				cli.XPS(job, args...)
			}
		},
	}

	var job string
	cmd.Flags().StringVarP(&job, "job", "j", "", "Show <JOB_ID> status")
	return cmd
}
