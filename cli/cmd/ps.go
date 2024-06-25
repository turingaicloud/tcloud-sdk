package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

func NewPSCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ps",
		Short: "Check submitted jobs' status",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			job, err := cmd.Flags().GetString("job")
			showK8SJob, _ := cmd.Flags().GetBool("showK8SJob")
			if  err == nil {
				cli.XPS(job, showK8SJob, args...)
			}else {
				cli.XPS("", showK8SJob, args...)
			}
		},
	}
	var showK8SJob bool
	var job string
	cmd.Flags().StringVarP(&job, "job", "j", "", "Show <JOB_ID> status")
	cmd.Flags().BoolVarP(&showK8SJob, "showK8SJob", "k", false, "To show Kubernetes task instead of Slurm job [ Experimental   ] (by default using Slurm as task engine)")
	return cmd
}
