package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

func NewSubmitCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	return &cobra.Command{
		Use:   "submit",
		Short: "Submit a job",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("tcloud submit CLI")
		},
	}
}
