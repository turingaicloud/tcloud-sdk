package tcloudcli

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/ssh"
)

type TcloudCli struct {
	userConfig	*UserConfig
	sess 		*Session
	prefix		string
}

func (tcloudcli *TcloudCli) NewSession() *Session {
	buffer, err := ioutil.ReadFile(tcloudcli.userConfig.authFile)
	if err != nil {
		fmt.Println("Failed to read authFile at %s", tcloudcli.userConfig.authFile)
		return nil
	}
	signer, _ := ssh.ParsePrivateKey(buffer)
	clientConfig := &ssh.ClientConfig{
		User: tcloudcli.userConfig.UserName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer)
		},
	}
	// TODO(SSHpath[0] to be removed when to one hop)
	client, err := ssh.Dial("tcp", tcloudcli.userConfig.SSHpath[0] + ":22", clientConfig)
	if err != nil {
		fmt.Println("Failed to dial: " + err.Error())
		return nil
	}
	session, err := client.NewSession()
	if err != nil {
		fmt.Println("Failed to create session: " + err.Error())
		return nil
	}
	modes := ssh.TerminalModes {
		ssh.ECHO:			0,
		ssh.TTY_OP_ISPEED:	14400,
		ssh.TTY_OP_OSPEED:	14400,
	}
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatal("Failed to request for pseudo terminal: ", err)
		return nil
	}
	return &session
}

func (tcloudcli *TcloudCli) NewPrefix() {
	if len(tcloudcli.userConfig.SSHpath) < 2 {
		tcloudcli.prefix = ""
	}
	var str string
	for _, s := range tcloudcli.userConfig.SSHpath[1:] {
		str = str + fmt.Sprintf("ssh -A -t %s@%s ", tcloudcli.userConfig.UserName, s)
	}
	tcloudcli.prefix = str
}

func NewTcloudCli(userConfig *UserConfig) *TcloudCli {
	tcloudcli := &TcloudCli{
		userConfig: userConfig,
	}
	if tcloudcli.sess = tcloudcli.NewSession(); tcloudcli.sess == nil {
		fmt.Println("Failed to start remote session")
		os.Exit(-1)
	}
	tcloudcli.NewPrefix()
	return tcloudcli
}

func (tcloudcli *TcloudCli) XBuild(args ...string) {
	var config TuxivConfig
	if workDir, err := config.ParseTuxivConf(args); err == true {
		fmt.Println("Parse tuxiv config file failed.")
		os.Exit(-1)
	}
	repoName := strings.Split(workDir, "/")[-1]
	if err = tcloudcli.UploadRepo(workDir); err == true {
		fmt.Println("Upload repository env failed")
		os.Exit(-1)
	}
	if err = tcloudcli.CondaRemove(config.Environment.Name); err == true {
		fmt.Println("Remove conda env failed")
		os.Exit(-1)
	}
	if err = tcloudcli.CondaCreate(repoName, config.Environment.Name); err == true {
		fmt.Println("Create conda env failed")
		os.Exit(-1)
	}
}

func (tcloudcli *TcloudCli) SendToCluster(src string) (string, bool) {
	if f, err := os.Stat(src); err != nil {
		fmt.Println("Failed to send to cluster. %s not exists.", src)
		return "", true
	}
	prefix = "-i"
	if mode := f.Mode(); mode.IsDir() {
		prefix = "-r -i"
	}

	// TODO(A bit wrong when transmit file. Not the same directory as src)
	dst := fmt.Sprintf("%s@%s:/home/%s", tcloudcli.userConfig.UserName, tcloudcli.userConfig.SSHpath[0], tcloudcli.userConfig.UserName)
	if len(tcloudcli.userConfig.SSHpath) < 2 {
		cmd := exec.Command("scp", prefix, tcloudcli.userConfig.authFile, src, dst)
	} else {
		str := ""
		for _, s := range tcloudcli.userConfig.SSHpath[:-1] {
			str = str + fmt.Sprintf("ssh -A %s@%s ", tcloudcli.userConfig.UserName, s)
		}
		str = str + "-W \%h:\%p"
		proxycmd := fmt.Sprintf("ProxyCommand=\"%s\"", str)
		cmd := exec.Command("scp", prefix, tcloudcli.userConfig.authFile, "-o", proxycmd, src, dst)
	}
	if _, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return dst, true
	}
	return dst, false
}

func (tcloudcli *TcloudCli) UploadRepo(workDir string) bool {
	// cmd := exec.Command("scp", "-r", "-i", tcloudcli.userConfig.authFile, "../tcloud_job", "ubuntu@18.162.45.250:/home/ubuntu")
	// out, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("cmd.Run() failed with %s\n", err)
	// 	return true
	// }
	// fmt.Printf("Upload file to TACC JUMP out:\n%s\n", string(out))
	// bash_command := "scp -r /home/ubuntu/tcloud_job TACC1:/home/ubuntu"
	// cmd = exec.Command("ssh", "-i", tcloudcli.userConfig.authFile, "ubuntu@18.162.45.250", bash_command)
	// out, err = cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("cmd.Run() failed with %s\n", err)
	// 	return true
	// }
	// fmt.Printf("Upload file to TACC1 out:\n%s\n", string(out))
	// bash_command = "scp -r /home/ubuntu/tcloud_job TACC2:/home/ubuntu"
	// cmd = exec.Command("ssh", "-i", tcloudcli.userConfig.authFile, "ubuntu@18.162.45.250", bash_command)
	// out, err = cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("cmd.Run() failed with %s\n", err)
	// 	return true
	// }
	// fmt.Printf("Upload file to TACC2 out:\n%s\n", string(out))
	dst, err := tcloudcli.SendToCluster(workDir)
	if err == true {
		fmt.Println("Failed to upload repo to ", dst)
		return true
	}
	fmt.Println("Successfully upload repo to ", dst)
	return false
}

func (tcloudcli *TcloudCli) CondaCreate(repoName string, envName string) bool {
	// bash_command := "ssh TACC1 /home/ubuntu/miniconda3/bin/conda env create -f /home/ubuntu/tcloud_job/configurations/conda.yaml"
	// cmd := exec.Command("ssh", "-i", tcloudcli.userConfig.authFile, "ubuntu@18.162.45.250", bash_command)
	// out, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("cmd.Run() failed with %s\n", err)
	// 	return true
	// }
	// fmt.Printf("conda create on TACC1 out:\n%s\n", string(out))
	// bash_command = "ssh TACC2 /home/ubuntu/miniconda3/bin/conda env create -f /home/ubuntu/tcloud_job/configurations/conda.yaml"
	// cmd = exec.Command("ssh", "-i", tcloudcli.userConfig.authFile, "ubuntu@18.162.45.250", bash_command)
	// out, err = cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("cmd.Run() failed with %s\n", err)
	// 	return true
	// }
	// fmt.Printf("conda create on TACC2 out:\n%s\n", string(out))
	// return false

	homeDir := fmt.Sprintf("/home/%s", tcloudcli.userConfig.UserName)
	condaBin := fmt.Sprintf("%s/miniconda3/bin/conda", homeDir)
	condaYaml := fmt.Sprintf("%s/%s/configurations/conda.yaml", homeDir, repoName)
	cmd := fmt.Sprintf("%s %s env create -f %s", tcloudcli.prefix, condaBin, condaYaml)
	if err := tcloudcli.sess.Run(cmd); err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Println("Environment %s created.", envName)
	return false
}
func (tcloudcli *TcloudCli) CondaRemove(envName string) bool {
	// bash_command := "ssh TACC1 /home/ubuntu/miniconda3/bin/conda remove -n " + name + " --all -y"
	// cmd := exec.Command("ssh", "-i", tcloudcli.userConfig.authFile, "ubuntu@18.162.45.250", bash_command)
	// out, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("cmd.Run() failed with %s\n", err)
	// 	return true
	// }
	// fmt.Printf("conda remove on TACC1 out:\n%s\n", string(out))
	// bash_command = "ssh TACC2 /home/ubuntu/miniconda3/bin/conda remove -n " + name + " --all -y"
	// cmd = exec.Command("ssh", "-i", tcloudcli.userConfig.authFile, "ubuntu@18.162.45.250", bash_command)
	// out, err = cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("cmd.Run() failed with %s\n", err)
	// 	return true
	// }
	// fmt.Printf("conda remove on TACC2 out:\n%s\n", string(out))
	// return false

	homeDir := fmt.Sprintf("/home/%s", tcloudcli.userConfig.UserName)
	condaBin := fmt.Sprintf("%s/miniconda3/bin/conda", homeDir)
	cmd := fmt.Sprintf("%s %s remove -n %s --all -y", tcloudcli.prefix, condaBin, envName)
	if err := tcloudcli.sess.Run(cmd); err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Println("Previous environment %s removed.", envName)
	return false
}
