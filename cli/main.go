package main

import (
	"fmt"
	"github.com/spf13/cobra"

	"os"
	"os/exec"
	"path/filepath"
	"tcloud-sdk/cli/cmd"
	"tcloud-sdk/cli/tcloudcli"
)

var VERSION = "0.0.3"

func main() {
	home := homeDIR()
	TcloudInit(home)
	// userConfig := tcloudcli.NewUserConfig(filepath.Join(home, ".tcloud", "user.json"), filepath.Join(home, ".tcloud", "TACC.pem"))
	userConfig := tcloudcli.NewUserConfig(filepath.Join(home, ".tcloud", ".userconfig"))
	clusterConfig := tcloudcli.NewClusterConfig(filepath.Join(home, ".tcloud", ".clusterconfig"))
	cli := tcloudcli.NewTcloudCli(userConfig, clusterConfig)
	tcloudCmd := newTcloudCommand(cli)
	if err := tcloudCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newTcloudCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	var tcloudCmd = &cobra.Command{
		Use:   "tcloud",
		Short: "TACC Command-line Interface v" + VERSION,
		// Version: VERSION,
	}
	tcloudCmd.AddCommand(cmd.NewSubmitCommand(cli))
	tcloudCmd.AddCommand(cmd.NewConfigCommand(cli))
	tcloudCmd.AddCommand(cmd.NewPSCommand(cli))
	tcloudCmd.AddCommand(cmd.NewInitCommand(cli))
	tcloudCmd.AddCommand(cmd.NewDownloadCommand(cli))
	tcloudCmd.AddCommand(cmd.NewAddCommand(cli))
	tcloudCmd.AddCommand(cmd.NewInstallCommand(cli))
	// tcloudCmd.AddCommand(cmd.NewLogCommand(cli))
	// tcloudCmd.AddCommand(cmd.NewCopyCommand(cli))
	tcloudCmd.AddCommand(cmd.NewDatasetCommand(cli))
	tcloudCmd.AddCommand(cmd.NewLSCommand(cli))

	var Verbose bool
	tcloudCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	return tcloudCmd
}

func homeDIR() string {
	return os.Getenv("HOME")
}

func TcloudInit(home string) bool {
	tcloudDIR := fmt.Sprintf("%s/.tcloud", home)
	init_cmd := exec.Command("mkdir", "-p", tcloudDIR)
	if _, err := init_cmd.CombinedOutput(); err != nil {
		fmt.Println("Failed to mkdir at ", tcloudDIR, err)
		return false
	}
	file1, err := os.Open(filepath.Join(home, ".tcloud", ".userconfig"))
	defer file1.Close()
	if err != nil && os.IsNotExist(err) {
		os.Create(filepath.Join(home, ".tcloud", ".userconfig"))
	}
	file2, err := os.Open(filepath.Join(home, ".tcloud", ".clusterconfig"))
	defer file2.Close()
	if err != nil && os.IsNotExist(err) {
		os.Create(filepath.Join(home, ".tcloud", ".clusterconfig"))
	}
	return true
}
