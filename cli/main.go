package main

import (
	"fmt"
	"github.com/spf13/cobra"

	"os"
	"path/filepath"
	"tcloud-sdk/cli/cmd"
	"tcloud-sdk/cli/tcloudcli"
)

func main() {
	home := homeDIR()
	userConfig := tcloudcli.NewUserConfig(filepath.Join(home, ".tcloud", "user.json"), filepath.Join(home, ".tcloud", "TACC.pem"))
	cli := tcloudcli.NewTcloudCli(userConfig)
	tcloudCmd := newTcloudCommand(cli)
	if err := tcloudCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newTcloudCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	var tcloudCmd = &cobra.Command{
		Use:   "tcloud",
		Short: "TACC Command-line Interface v0.0.1",
	}
	tcloudCmd.AddCommand(cmd.NewBuildCommand(cli))
	tcloudCmd.AddCommand(cmd.NewSubmitCommand(cli))
	return tcloudCmd
}

func homeDIR() string {
	return os.Getenv("HOME")
}
