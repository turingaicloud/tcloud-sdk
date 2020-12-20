package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"

	"os"
	"os/exec"
	"path/filepath"
	"tcloud-sdk/cli/cmd"
	"tcloud-sdk/cli/tcloudcli"
)

var VERSION = "0.1.0"

func main() {
	home := homeDIR()
	TcloudInit(home)
	userConfig := tcloudcli.NewUserConfig(filepath.Join(home, ".tcloud", ".userconfig"))
	clusterConfig := tcloudcli.NewClusterConfig(filepath.Join(home, ".tcloud", ".clusterconfig"))
	cli := tcloudcli.NewTcloudCli(userConfig, clusterConfig)
	tcloudCmd := newTcloudCommand(cli)
	if err := tcloudCmd.Execute(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

func newTcloudCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	var tcloudCmd = &cobra.Command{
		Use:   "tcloud",
		Short: "TACC Command-line Interface v" + VERSION,
	}
	tcloudCmd.AddCommand(cmd.NewSubmitCommand(cli))
	tcloudCmd.AddCommand(cmd.NewConfigCommand(cli))
	tcloudCmd.AddCommand(cmd.NewPSCommand(cli))
	tcloudCmd.AddCommand(cmd.NewInitCommand(cli))
	tcloudCmd.AddCommand(cmd.NewDownloadCommand(cli))
	tcloudCmd.AddCommand(cmd.NewAddCommand(cli))
	tcloudCmd.AddCommand(cmd.NewInstallCommand(cli))
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
	log.SetPrefix("[Tcloud]")
	log.SetFlags(log.Ldate | log.Lshortfile)

	tcloudDIR := fmt.Sprintf("%s/.tcloud", home)
	init_cmd := exec.Command("mkdir", "-p", tcloudDIR)
	if _, err := init_cmd.CombinedOutput(); err != nil {
		log.Println("Failed to obtain tcloud metadata. Error message:", err.Error())
		return false
	}
	file1, err := os.Open(filepath.Join(home, ".tcloud", ".userconfig"))
	if err != nil && os.IsNotExist(err) {
		os.Create(filepath.Join(home, ".tcloud", ".userconfig"))
	}
	file2, err := os.Open(filepath.Join(home, ".tcloud", ".clusterconfig"))
	if err != nil && os.IsNotExist(err) {
		os.Create(filepath.Join(home, ".tcloud", ".clusterconfig"))
	}

	defer file1.Close()
	defer file2.Close()

	return true
}
