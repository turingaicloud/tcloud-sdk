package cmd

import (
	"tcloud-sdk/cli/tcloudcli"

	"github.com/spf13/cobra"
)

func NewCatCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cat",
		Short: "Concatenate FILE(s) to standard output.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cli.XCat(args...)
		},
	}

	var job string
	cmd.Flags().StringVarP(&job, "job", "j", "", "Show <JOB_ID> status")
	return cmd
}
