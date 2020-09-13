package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
	"tcloud-sdk/cli/tcloudcli"
)

func NewSubmitCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	return &cobra.Command{
		Use:   "submit",
		Short: "Submit a job",
		Run: func(cmd *cobra.Command, args []string) {
			// fmt.Println("tcloud submit CLI")
			// run TACC1:/home/ubuntu/mnist_demo/run.slurm
			err := SubmitJob()
			if err {
				fmt.Println("Fail to submit job.")
				log.Fatal(err)
			}
		},
	}
}

func SubmitJob() bool {
	bash_command := "ssh ubuntu@TACC1 \"sbatch /home/ubuntu/mnist_demo/run.slurm\""
	cmd := exec.Command("ssh", "-i", "/Users/xcwan/Downloads/TACC.pem", "ubuntu@18.162.45.250", bash_command)
	fmt.Println(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Printf("Job submitted\n", string(out))
	return false
}
