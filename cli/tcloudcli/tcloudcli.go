package tcloudcli

import (
	"fmt"
	"io"
	"io/ioutil"
	// "log"
	"bytes"
	"net"
	"os"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

type SSHTunnel struct {
	sess *ssh.Session
	in   *io.Reader
	out  *io.Writer
	err  *io.Writer
}

type TcloudCli struct {
	userConfig *UserConfig
	tunnel     *SSHTunnel
	prefix     string
}

func (tcloudcli *TcloudCli) NewSession() *SSHTunnel {
	var tunnel *SSHTunnel
	buffer, err := ioutil.ReadFile(tcloudcli.userConfig.authFile)
	if err != nil {
		fmt.Println("Failed to read authFile at %s", tcloudcli.userConfig.authFile)
		return nil
	}
	signer, _ := ssh.ParsePrivateKey(buffer)
	clientConfig := &ssh.ClientConfig{
		User: tcloudcli.userConfig.UserName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	// TODO(SSHpath[0] to be removed when to one hop)
	client, err := ssh.Dial("tcp", tcloudcli.userConfig.SSHpath[0]+":22", clientConfig)
	if err != nil {
		fmt.Println("Failed to dial: " + err.Error())
		return nil
	}
	session, err := client.NewSession()
	if err != nil {
		fmt.Println("Failed to create session: " + err.Error())
		return nil
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		fmt.Println("Failed to request for pseudo terminal: ", err)
		return nil
	}
	// Print stdout
	sess_stdin, err := session.StdinPipe()
	sess_stdout, err := session.StdoutPipe()
	sess_stderr, err := session.StderrPipe()
	tunnel.sess = session
	tunnel.in = &sess_stdin
	tunnel.out = &sess_stdout
	tunnel.err = &sess_stderr

	// Start remote shell
	err = tunnel.sess.Shell()
	if err != nil {
		fmt.Println("Failed to start remote shell: ", err)
		return nil
	}
	return tunnel
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
	if tcloudcli.tunnel = tcloudcli.NewSession(); tcloudcli.tunnel == nil {
		fmt.Println("Failed to start remote session")
		os.Exit(-1)
	}
	tcloudcli.NewPrefix()
	return tcloudcli
}

func (tcloudcli *TcloudCli) XBuild(args ...string) {
	var config TuxivConfig
	localWorkDir, repoName, err := config.ParseTuxivConf(tcloudcli, args)
	if err == true {
		fmt.Println("Parse tuxiv config file failed.")
		os.Exit(-1)
	}
	if err = tcloudcli.UploadRepo(localWorkDir); err == true {
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
	f, err := os.Stat(src)
	if err != nil {
		fmt.Println("Failed to send to cluster. %s not exists.", src)
		return "", true
	}
	prefix := ""
	if mode := f.Mode(); mode.IsDir() {
		prefix = "-r"
	}

	// TODO(A bit wrong when transmit file. Not the same directory as src)
	dst := tcloudcli.userConfig.SSHpath[len(tcloudcli.userConfig.SSHpath)-1]
	dst = fmt.Sprintf("%s@%s:/home/%s", tcloudcli.userConfig.UserName, dst, tcloudcli.userConfig.UserName)
	cmd := exec.Command("")
	if len(tcloudcli.userConfig.SSHpath) < 2 {
		cmd = exec.Command("scp", "-i", tcloudcli.userConfig.authFile, prefix, src, dst)
	} else {
		str := ""
		sshpath := tcloudcli.userConfig.SSHpath
		for _, s := range sshpath[:len(sshpath)-1] {
			str = str + fmt.Sprintf(`ssh -i %s -o StrictHostKeyChecking=no -A -t %s@%s `, tcloudcli.userConfig.authFile, tcloudcli.userConfig.UserName, s)
		}
		// str = str + "-W %h:%p"
		proxycmd := fmt.Sprintf(`ProxyCommand="%s"`, str)
		cmd = exec.Command("scp", "-i", tcloudcli.userConfig.authFile, "-o", proxycmd, ` -o StrictHostKeyChecking=no`, prefix, src, dst)
	}
	fmt.Println(cmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		// if exitError, ok := err.(*exec.ExitError); ok {
		// 	fmt.Println("Failed to run cmd in SendToCluster ", exitError.ExitCode())
		// 	return dst, true
		// }
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		fmt.Println("Failed to run cmd in SendToCluster ", err)
		return dst, true
	}
	return dst, false
}

func (tcloudcli *TcloudCli) UploadRepo(localWorkDir string) bool {
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
	dst, err := tcloudcli.SendToCluster(localWorkDir)
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
	cmd := fmt.Sprintf("%s %s env create -f %s\n", tcloudcli.prefix, condaBin, condaYaml)
	if err := tcloudcli.tunnel.sess.Run(cmd); err != nil {
		fmt.Println("Failed to run cmd in CondaCreate ", err)
		return true
	}
	write(tcloudcli.tunnel.in, cmd)
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
		fmt.Println("Failed to run cmd in CondaRemove ", err)
		return true
	}
	fmt.Println("Previous environment \"", envName, "\" removed.")
	return false
}
