package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"tcloud-sdk/cli/tcloudcli"
)

func NewConfigCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure user's account in tcloud CLI",
		Run: func(cmd *cobra.Command, args []string) {
			if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
				fmt.Println("user.name:\t", cli.UserConfig("username")[0])
				fmt.Println("user.authfile:\t", cli.UserConfig("authfile")[0])
			}
			userName, _ := cmd.Flags().GetString("username")
			authFile, _ := cmd.Flags().GetString("file")

			// Modify username and authfile, but maintain SSH path and Dir. And update userconfig file
			var config tcloudcli.UserConfig
			file, err := os.Create(cli.UserConfig("path")[0])
			defer file.Close()
			if err != nil {
				fmt.Println("Failed to open ", cli.UserConfig("path")[0])
				return
			}
			config.UserName = userName
			config.AuthFile = authFile
			config.SSHpath = cli.UserConfig("sshpath")
			// config.Dir = cli.ClusterConfig("dir")
			// config.path = cli.UserConfig("path")[0]

			encoder := json.NewEncoder(file)
			if err := encoder.Encode(config); err != nil {
				fmt.Println("Failed to encode struct UserConfig.", err.Error())
				return
			}
		},
	}

	var userName string
	var authFile string
	cmd.Flags().StringVarP(&userName, "username", "u", cli.UserConfig("username")[0], "Configure TACC username")
	cmd.Flags().StringVarP(&authFile, "file", "f", cli.UserConfig("authfile")[0], "Configure TACC authentication file path")

	return cmd
}
