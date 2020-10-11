package tcloudcli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	shellquote "github.com/gonuts/go-shellquote"
	"golang.org/x/crypto/ssh"
)

type TcloudCli struct {
	userConfig *UserConfig
	// sess       *ssh.Session
	prefix string
}

func (tcloudcli *TcloudCli) UserConfig(option string) []string {
	var s []string
	switch strings.ToLower(option) {
	case "username":
		return append(s, tcloudcli.userConfig.UserName)
	case "authfile":
		return append(s, tcloudcli.userConfig.AuthFile)
	case "sshpath":
		return tcloudcli.userConfig.SSHpath
	case "path":
		return append(s, tcloudcli.userConfig.path)
	default:
		fmt.Println("No options found in userconfig")
		return s
	}
}

func (tcloudcli *TcloudCli) NewSession() *ssh.Session {
	buffer, err := ioutil.ReadFile(tcloudcli.userConfig.AuthFile)
	if err != nil {
		fmt.Println("Failed to read AuthFile at %s", tcloudcli.userConfig.AuthFile)
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
	} else {
		var str string
		for _, s := range tcloudcli.userConfig.SSHpath[1:] {
			str = str + fmt.Sprintf("ssh -A -t %s@%s ", tcloudcli.userConfig.UserName, s)
		}
		tcloudcli.prefix = str
	}

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
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	if err := sess.Run(cmd); err != nil {
		fmt.Println("Failed to run cmd \"", cmd, "\"", err)
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

func (tcloudcli *TcloudCli) SendRepoToCluster(repoName string, src string) (string, bool) {
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
	cmd := exec.Command("scp", prefix, "-i", tcloudcli.userConfig.AuthFile, src, dst)
	if _, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("Failed to run cmd in SendRepoToCluster ", err)
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

// SCP from SSHPath[0] to localhost
func (tcloudcli *TcloudCli) RecvFromCluster(src string, dst string, IsDir bool) bool {
	srcIP := tcloudcli.userConfig.SSHpath[0]
	srcPath := fmt.Sprintf("%s@%s:%s", tcloudcli.userConfig.UserName, srcIP, src)
	dstPath := fmt.Sprintf("%s", dst)

	var cmd *exec.Cmd
	if IsDir {
		cmd = exec.Command("scp", "-r", "-i", tcloudcli.userConfig.AuthFile, srcPath, dstPath)
	} else {
		cmd = exec.Command("scp", "-i", tcloudcli.userConfig.AuthFile, srcPath, dstPath)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to run cmd in RecvFromCluster ", err.Error(), stderr.String())
		return true
	}
	return false
}

func (tcloudcli *TcloudCli) BuildEnv(args ...string) {
	var config TuxivConfig
	localWorkDir, repoName, err := config.ParseTuxivConf(tcloudcli, args)
	randString := RandString(16)
	if err == true {
		fmt.Println("Parse tuxiv config file failed.")
		os.Exit(-1)
	}
	if err = tcloudcli.UploadRepo(repoName, localWorkDir); err == true {
		fmt.Println("Upload repository env failed")
		os.Exit(-1)
	}
	if err = tcloudcli.CondaRemove(config.Environment.Name, randString); err == true {
		fmt.Println("Remove conda env failed")
		os.Exit(-1)
	}
	if err = tcloudcli.CondaCreate(repoName, config.Environment.Name, randString); err == true {
		fmt.Println("Create conda env failed")
		os.Exit(-1)
	}
}
func (tcloudcli *TcloudCli) UploadRepo(repoName string, localWorkDir string) bool {
	dst, err := tcloudcli.SendRepoToCluster(repoName, localWorkDir)
	if err == true {
		fmt.Println("Failed to upload repo to ", dst)
		return true
	}
	fmt.Println("Successfully upload repo to ", dst)
	return false
}
func (tcloudcli *TcloudCli) CondaCreate(repoName string, envName string, randString string) bool {
	homeDir := fmt.Sprintf("/home/%s", tcloudcli.userConfig.UserName)
	condaBin := fmt.Sprintf("%s/miniconda3/bin/conda", homeDir)
	condaYaml := fmt.Sprintf("%s/%s/configurations/conda.yaml", homeDir, repoName)
	cmd := fmt.Sprintf("%s %s env create -f %s -n %s\n", tcloudcli.prefix, condaBin, condaYaml, envName+randString)
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		fmt.Println("Failed to run cmd in CondaCreate ", err)
		return true
	}

	fmt.Println("Environment \"", envName, "\" created.")
	return false
}
func (tcloudcli *TcloudCli) CondaRemove(envName string, randString string) bool {
	homeDir := fmt.Sprintf("/home/%s", tcloudcli.userConfig.UserName)
	condaBin := fmt.Sprintf("%s/miniconda3/bin/conda", homeDir)
	cmd := fmt.Sprintf("%s %s remove -n %s --all -y", tcloudcli.prefix, condaBin, envName+randString)
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		fmt.Println("Failed to run cmd in CondaRemove ", err)
		return true
	}
	fmt.Println("Previous environment \"", envName, "\" removed.")
	return false
}

func RandString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6
		letterIdxMask = 1<<letterIdxBits - 1
		letterIdxMax  = 63 / letterIdxBits
	)
	var src = rand.NewSource(time.Now().UnixNano())
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}
func (tcloudcli *TcloudCli) XSubmit(args ...string) bool {
	tcloudcli.BuildEnv(args...)
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

func (tcloudcli *TcloudCli) XPS(job string, args ...string) bool {
	var cmd string
	if job == "" {
		cmd = fmt.Sprintf("%s squeue", tcloudcli.prefix)
	} else {
		cmd = fmt.Sprintf("%s squeue -j %s", tcloudcli.prefix, job)
	}
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		fmt.Println("Failed to run cmd in tcloud ps.")
		return true
	}
	return false
}

// TODO(Just a receive file prototype, dstPath TODEFINE, config TODEFINE)
func (tcloudcli *TcloudCli) XInit(args ...string) bool {
	// TODO(config file path)
	if len(args) == 1 {
		// Remote receive config file
		src := fmt.Sprintf("/home/%s/%s/main.go", tcloudcli.userConfig.UserName, args[0])
		dst := fmt.Sprintf("%s", filepath.Join(os.Getenv("HOME"), ".tcloud"))
		IsDir := false

		cmd := fmt.Sprintf("scp %s@%s:%s %s", tcloudcli.userConfig.UserName, tcloudcli.userConfig.SSHpath[1], src, src)
		if err := tcloudcli.RemoteExecCmd(cmd); err == true {
			fmt.Println("Failed to receive file at Staging Node: ", err)
			return true
		}

		if err := tcloudcli.RecvFromCluster(src, dst, IsDir); err == true {
			fmt.Println("Failed to receive file at localhost.")
			return true
		}
		// TODO(Parse config file and update shell)
		fmt.Println("User's file configured.")
		return false
	}
	fmt.Println("Failed to parse args.")
	return true
}

func (tcloudcli *TcloudCli) XDownload(args ...string) bool {
	cmd := exec.Command("wget", args[0])
	if _, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("Failed to wget", args[0])
		return true
	}
	return false
}
func (tcloudcli *TcloudCli) XAdd(args ...string) bool {
	// Add new dependency to tuxiv.conf
	var config TuxivConfig
	err := config.AddDepTuxivFile(tcloudcli, args)
	if err == true {
		fmt.Println("Add dependency to tuxiv config file failed.")
		os.Exit(-1)
	}
	return false
}
func (tcloudcli *TcloudCli) XInstall(args ...string) bool {
	var config TuxivConfig
	_, _, err := config.ParseTuxivConf(tcloudcli, args)
	if err == true {
		fmt.Println("Parse tuxiv config file failed.")
		os.Exit(-1)
	}
	condaYaml := fmt.Sprintf("./configurations/conda.yaml")
	removeCmd := exec.Command("conda", "env", "remove", "-n", config.Environment.Name)
	if out, err := removeCmd.CombinedOutput(); err != nil {
		fmt.Println("Failed to create local environment err: ", err)
		return true
	} else {
		fmt.Printf("%s\n", string(out))
	}
	createCmd := exec.Command("conda", "env", "create", "-f", condaYaml)
	if out, err := createCmd.CombinedOutput(); err != nil {
		fmt.Println("Failed to create local environment err: ", err)
		return true
	} else {
		fmt.Printf("%s\n", string(out))
	}
	fmt.Println("Environment \"", config.Environment.Name, "\" created locally.")
	return false
}
