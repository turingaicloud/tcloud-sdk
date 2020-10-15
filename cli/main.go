package main

import (
	"fmt"
	"github.com/spf13/cobra"

	"os"
	"path/filepath"
	"tcloud-sdk/cli/cmd"
	"tcloud-sdk/cli/tcloudcli"
)

var VERSION = "0.0.2"

func main() {
	home := homeDIR()
	// userConfig := tcloudcli.NewUserConfig(filepath.Join(home, ".tcloud", "user.json"), filepath.Join(home, ".tcloud", "TACC.pem"))
	userConfig := tcloudcli.NewUserConfig(filepath.Join(home, ".tcloud", ".userconfig"))
	cli := tcloudcli.NewTcloudCli(userConfig)
	tcloudCmd := newTcloudCommand(cli)
	if err := tcloudCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newTcloudCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	var tcloudCmd = &cobra.Command{
		Use:     "tcloud",
		Short:   "TACC Command-line Interface v" + VERSION,
		Version: VERSION,
	}
	tcloudCmd.AddCommand(cmd.NewSubmitCommand(cli))
	tcloudCmd.AddCommand(cmd.NewConfigCommand(cli))
	tcloudCmd.AddCommand(cmd.NewPSCommand(cli))
	tcloudCmd.AddCommand(cmd.NewInitCommand(cli))
	tcloudCmd.AddCommand(cmd.NewDownloadCommand(cli))
	tcloudCmd.AddCommand(cmd.NewAddCommand(cli))
	tcloudCmd.AddCommand(cmd.NewInstallCommand(cli))

	var Verbose bool
	tcloudCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	return tcloudCmd
}

func homeDIR() string {
	return os.Getenv("HOME")
}
