package cmd

import (
	"tcloud-sdk/cli/tcloudcli"

	"github.com/spf13/cobra"
)

func NewENVLSCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "List environment",
	}

	cmd.AddCommand(sub(cli))
	return cmd
}

func sub(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List directory contents",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			IsLong, _ := cmd.Flags().GetBool("long")
			IsReverse, _ := cmd.Flags().GetBool("reverse")
			IsAll, _ := cmd.Flags().GetBool("all")
			cli.XENVLS(IsLong, IsReverse, IsAll, args...)
		},
	}

	var IsLong, IsReverse, IsAll bool
	cmd.Flags().BoolVarP(&IsLong, "long", "l", false, "List in long format.")
	cmd.Flags().BoolVarP(&IsReverse, "reverse", "r", false, "Reverse the order of the sort.")
	cmd.Flags().BoolVarP(&IsAll, "all", "a", false, "Include directory entries whose names begin with a dot (.).")
	return cmd
}
