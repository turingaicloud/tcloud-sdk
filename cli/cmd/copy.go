package cmd

import (
	"github.com/spf13/cobra"
	"tcloud-sdk/cli/tcloudcli"
)

func NewCopyCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy file/directory to user's directory",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			IsDir, _ := cmd.Flags().GetBool("recursive")
			cli.XCopy(IsDir, args...)
		},
	}
	var IsDir bool
	cmd.Flags().BoolVarP(&IsDir, "recursive", "r", false, "recursively copy")

	return cmd
}
