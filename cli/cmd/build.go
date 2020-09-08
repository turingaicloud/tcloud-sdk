package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Setup conda environment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tcloud build CLI")
	},
}
