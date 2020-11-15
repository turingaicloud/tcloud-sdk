package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

func NewLSCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List directory contents",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			IsLong, _ := cmd.Flags().GetBool("long")
			IsReverse, _ := cmd.Flags().GetBool("reverse")
			IsAll, _ := cmd.Flags().GetBool("all")
			cli.XLS(IsLong, IsReverse, IsAll, args...)
		},
	}
	var IsLong, IsReverse, IsAll bool
	cmd.Flags().BoolVarP(&IsLong, "long", "l", false, "List in long format.")
	cmd.Flags().BoolVarP(&IsReverse, "reverse", "r", false, "Reverse the order of the sort.")
	cmd.Flags().BoolVarP(&IsAll, "all", "a", false, "Include directory entries whose names begin with a dot (.).")
	return cmd
}
