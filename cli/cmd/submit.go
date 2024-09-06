package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

func NewSubmitCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit",
		Short: "Submit a job to TACC",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			submitToK8s, _ := cmd.Flags().GetBool("submitToK8s")
			cli.XSubmit(submitToK8s, args...)
		},
	}
	var submitToK8s bool
	cmd.Flags().BoolVarP(&submitToK8s, "submitToK8s", "k", false, "Task submitted to Kubernetes? [ Experimental ] (by default using Slurm as task engine)")
	return cmd
}
