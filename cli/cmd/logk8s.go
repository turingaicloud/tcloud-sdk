package cmd

import (
	"tcloud-sdk/cli/tcloudcli"
	"log"
	"github.com/spf13/cobra"
)

func NewLogK8SCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "dump logs for Kubernetes jobs",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			job, err := cmd.Flags().GetString("job")
			if err != nil || job == "" {
				log.Println("Error: parameter job (--job) should not be blank.")
				return
			}else {
				cli.XLogK8S(job, args...)
			}
		},
	}

	var job string
	cmd.Flags().StringVarP(&job, "job", "j", "", "Show <K8S_JOB_Name> status")
	return cmd
}
