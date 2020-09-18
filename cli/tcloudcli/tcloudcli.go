package tcloudcli

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	shellquote "github.com/gonuts/go-shellquote"
	"golang.org/x/crypto/ssh"
)

type TcloudCli struct {
	userConfig *UserConfig
	// sess       *ssh.Session
	prefix string
}

func (tcloudcli *TcloudCli) NewSession() *ssh.Session {
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
	return session
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
	tcloudcli.NewPrefix()
	return tcloudcli
}

func (tcloudcli *TcloudCli) RemoteExecCmd(cmd string) bool {
	sess := tcloudcli.NewSession()
	if sess == nil {
		fmt.Println("Failed to create remote session")
		os.Exit(-1)
	}
	w, err := sess.StdinPipe()
	if err != nil {
		fmt.Println("Failed to create StdinPipe", err)
		return true
	}

	if err := sess.Run(cmd); err != nil {
		fmt.Println("Failed to run cmd in SendToCluster ", err)
		w.Close()
		return true
	}
	defer sess.Close()

	errors := make(chan error)
	go func() {
		errors <- sess.Wait()
	}()
	fmt.Fprint(w, "\x00")
	w.Close()
	return false
}

func (tcloudcli *TcloudCli) SendToCluster(repoName string, src string) (string, bool) {
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
	dst := tcloudcli.userConfig.SSHpath[0]
	dst = fmt.Sprintf("%s@%s:/home/%s/", tcloudcli.userConfig.UserName, dst, tcloudcli.userConfig.UserName)
	cmd := exec.Command("scp", prefix, "-i", tcloudcli.userConfig.authFile, src, dst)
	if _, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("Failed to run cmd in SendToCluster ", err)
		return dst, true
	}
	if len(tcloudcli.userConfig.SSHpath) < 2 {
		return dst, false
	} else if len(tcloudcli.userConfig.SSHpath) == 2 {
		src := fmt.Sprintf("/home/%s/%s", tcloudcli.userConfig.UserName, repoName)
		dst := tcloudcli.userConfig.SSHpath[1]
		dst = fmt.Sprintf("%s@%s:/home/%s/", tcloudcli.userConfig.UserName, dst, tcloudcli.userConfig.UserName)
		cmd := shellquote.Join("scp", prefix, src, dst)
		if err := tcloudcli.RemoteExecCmd(cmd); err == true {
			fmt.Println("Failed to send repo from Host:", tcloudcli.userConfig.SSHpath[0], " to Host:", tcloudcli.userConfig.SSHpath[1])
			return dst, true
		}
		return dst, false
	} else {
		fmt.Println("Not support multi-hop send")
		return dst, true
	}
}

func (tcloudcli *TcloudCli) XBuild(args ...string) {
	var config TuxivConfig
	localWorkDir, repoName, err := config.ParseTuxivConf(tcloudcli, args)
	if err == true {
		fmt.Println("Parse tuxiv config file failed.")
		os.Exit(-1)
	}
	if err = tcloudcli.UploadRepo(repoName, localWorkDir); err == true {
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
func (tcloudcli *TcloudCli) UploadRepo(repoName string, localWorkDir string) bool {
	dst, err := tcloudcli.SendToCluster(repoName, localWorkDir)
	if err == true {
		fmt.Println("Failed to upload repo to ", dst)
		return true
	}
	fmt.Println("Successfully upload repo to ", dst)
	return false
}
func (tcloudcli *TcloudCli) CondaCreate(repoName string, envName string) bool {
	homeDir := fmt.Sprintf("/home/%s", tcloudcli.userConfig.UserName)
	condaBin := fmt.Sprintf("%s/miniconda3/bin/conda", homeDir)
	condaYaml := fmt.Sprintf("%s/%s/configurations/conda.yaml", homeDir, repoName)
	cmd := fmt.Sprintf("%s %s env create -f %s\n", tcloudcli.prefix, condaBin, condaYaml)
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		fmt.Println("Failed to run cmd in CondaCreate ", err)
		return true
	}

	fmt.Println("Environment \"", envName, "\" created.")
	return false
}
func (tcloudcli *TcloudCli) CondaRemove(envName string) bool {
	homeDir := fmt.Sprintf("/home/%s", tcloudcli.userConfig.UserName)
	condaBin := fmt.Sprintf("%s/miniconda3/bin/conda", homeDir)
	cmd := fmt.Sprintf("%s %s remove -n %s --all -y", tcloudcli.prefix, condaBin, envName)
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		fmt.Println("Failed to run cmd in CondaRemove ", err)
		return true
	}
	fmt.Println("Previous environment \"", envName, "\" removed.")
	return false
}

func (tcloudcli *TcloudCli) XSubmit(args ...string) bool {
	repoName := ""
	if len(args) < 1 {
		localWorkDir, _ := filepath.Abs(".")
		dirlist := strings.Split(localWorkDir, "/")
		repoName = dirlist[len(dirlist)-1]
	} else {
		repoName = args[0]
	}
	homeDir := fmt.Sprintf("/home/%s/%s", tcloudcli.userConfig.UserName, repoName)
	cmd := fmt.Sprintf("%s sbatch %s/configurations/run.slurm", tcloudcli.prefix, homeDir)
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		fmt.Println("Failed to run cmd in tcloud submit: ", err)
		return true
	}
	fmt.Println("Job", repoName, "submitted.")
	return false
}
