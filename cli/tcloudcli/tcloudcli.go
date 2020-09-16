package tcloudcli

import (
	"bufio"
	"log"
	"os"
	"os/exec"
)

type TcloudCli struct {
	userConfig *UserConfig
}

func NewTcloudCli(userConfig *UserConfig) *TcloudCli {
	tcloudcli := &TcloudCli{
		userConfig: userConfig,
	}
	return tcloudcli
}

// Build command
func (tcloudcli *TcloudCli) XBuild(args ...string) {
	var config TuxivConfig
	if err := config.ParseTuxivConf(args); err == true {
		fmt.Println("Parse tuxiv config file failed.")
		os.Exit(-1)
	}
	if err = tcloudcli.UploadFile(); err == true {
		fmt.Println("Upload file env failed")
		os.Exit(-1)
	}
	if err = tcloudcli.CondaRemove(setting.Environment.Name); err == true {
		fmt.Println("Remove conda env failed")
		os.Exit(-1)
	}
	if err = tcloudcli.CondaCreate(); err == true {
		fmt.Println("Create conda env failed")
		os.Exit(-1)
	}
}

func (tcloudcli *TcloudCli) UploadFile() bool {
	cmd := exec.Command("scp", "-r", "-i", tcloudcli.userConfig.authFile, "../tcloud_job", "ubuntu@18.162.45.250:/home/ubuntu")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Printf("Upload file to TACC JUMP out:\n%s\n", string(out))
	bash_command := "scp -r /home/ubuntu/tcloud_job TACC1:/home/ubuntu"
	cmd = exec.Command("ssh", "-i", tcloudcli.userConfig.authFile, "ubuntu@18.162.45.250", bash_command)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Printf("Upload file to TACC1 out:\n%s\n", string(out))
	bash_command = "scp -r /home/ubuntu/tcloud_job TACC2:/home/ubuntu"
	cmd = exec.Command("ssh", "-i", tcloudcli.userConfig.authFile, "ubuntu@18.162.45.250", bash_command)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Printf("Upload file to TACC2 out:\n%s\n", string(out))
	return false
}
func (tcloudcli *TcloudCli) CondaCreate() bool {
	bash_command := "ssh TACC1 /home/ubuntu/miniconda3/bin/conda env create -f /home/ubuntu/tcloud_job/configurations/conda.yaml"
	cmd := exec.Command("ssh", "-i", tcloudcli.userConfig.authFile, "ubuntu@18.162.45.250", bash_command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Printf("conda create on TACC1 out:\n%s\n", string(out))
	bash_command = "ssh TACC2 /home/ubuntu/miniconda3/bin/conda env create -f /home/ubuntu/tcloud_job/configurations/conda.yaml"
	cmd = exec.Command("ssh", "-i", tcloudcli.userConfig.authFile, "ubuntu@18.162.45.250", bash_command)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Printf("conda create on TACC2 out:\n%s\n", string(out))
	return false
}
func (tcloudcli *TcloudCli) CondaRemove(name string) bool {
	bash_command := "ssh TACC1 /home/ubuntu/miniconda3/bin/conda remove -n " + name + " --all -y"
	cmd := exec.Command("ssh", "-i", tcloudcli.userConfig.authFile, "ubuntu@18.162.45.250", bash_command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Printf("conda remove on TACC1 out:\n%s\n", string(out))
	bash_command = "ssh TACC2 /home/ubuntu/miniconda3/bin/conda remove -n " + name + " --all -y"
	cmd = exec.Command("ssh", "-i", tcloudcli.userConfig.authFile, "ubuntu@18.162.45.250", bash_command)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Printf("conda remove on TACC2 out:\n%s\n", string(out))
	return false
}
