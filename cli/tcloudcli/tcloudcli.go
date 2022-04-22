package tcloudcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

var DEFAULT_CLUSTERCONFIG_PATH = "/mnt/data/.clusterconfig"
var CityNetAPI = "http://gw.turing.ust.hk:9002/datasets"

type TcloudCli struct {
	userConfig    *UserConfig
	clusterConfig *ClusterConfig
	// globalenv     *TACCGlobalEnv
	prefix string
}

type Dataset struct {
	Name        string   `json: "name"`
	ID          string   `json: "_id"`
	CreateTime  string   `json: "create_time"`
	Files       []string `json: "files"`
	Labels      []string `json: "labels"`
	Description string   `json: "description"`
	Categories  string   `json: "categories"`
	Path        string   `json: "path"`
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
	case "port":
		return append(s, tcloudcli.userConfig.Port)
	case "path":
		return append(s, tcloudcli.userConfig.path)
	default:
		log.Println("No option found in userconfig.")
		return s
	}
}
func (tcloudcli *TcloudCli) ClusterConfig(option string) []string {
	var s []string
	switch strings.ToLower(option) {
	case "tcloudversion":
		return append(s, tcloudcli.clusterConfig.TcloudVersion)
	case "dirs":
		return append(s, tcloudcli.clusterConfig.Dirs["workdir"], tcloudcli.clusterConfig.Dirs["userdir"])
	case "homedir":
		return append(s, tcloudcli.clusterConfig.HomeDir)
	case "datasetdir":
		return append(s, tcloudcli.clusterConfig.DatasetDir)
	case "conda":
		return append(s, tcloudcli.clusterConfig.Conda)
	default:
		log.Println("No option found in clusterconfig.")
		return s
	}
}

func (tcloudcli *TcloudCli) NewSession() *ssh.Session {
	buffer, err := ioutil.ReadFile(tcloudcli.userConfig.AuthFile)
	if err != nil {
		log.Println("Failed to read AuthFile at ", tcloudcli.userConfig.AuthFile)
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
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", tcloudcli.userConfig.SSHpath[0], tcloudcli.userConfig.Port), clientConfig)
	if err != nil {
		log.Println("Failed to dial: " + err.Error())
		return nil
	}
	session, err := client.NewSession()
	if err != nil {
		log.Println("Failed to create session: " + err.Error())
		return nil
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Println("Failed to request for pseudo terminal: " + err.Error())
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

func NewTcloudCli(userConfig *UserConfig, clusterConfig *ClusterConfig) *TcloudCli {
	tcloudcli := &TcloudCli{
		userConfig:    userConfig,
		clusterConfig: clusterConfig,
	}
	tcloudcli.NewPrefix()
	return tcloudcli
}

func (tcloudcli *TcloudCli) RemoteExecCmd(cmd string) bool {
	sess := tcloudcli.NewSession()
	if sess == nil {
		log.Println("Failed to create remote session")
		os.Exit(-1)
	}
	w, err := sess.StdinPipe()
	if err != nil {
		log.Println("Failed to create StdinPipe", err)
		return true
	}
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	if err := sess.Run(cmd); err != nil {
		log.Println("Failed to run cmd \"", cmd, "\"", err)
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

func (tcloudcli *TcloudCli) RemoteExecCmdOutput(cmd string) ([]byte, bool) {
	sess := tcloudcli.NewSession()
	if sess == nil {
		log.Println("Failed to create remote session")
		os.Exit(-1)
	}
	w, err := sess.StdinPipe()
	if err != nil {
		log.Println("Failed to create StdinPipe", err)
		return nil, true
	}
	var b bytes.Buffer

	sess.Stdout = &b
	sess.Stderr = os.Stderr

	if err := sess.Run(cmd); err != nil {
		log.Println("Failed to run cmd \"", cmd, "\"", err)
		w.Close()
		return nil, true
	}
	defer sess.Close()

	errors := make(chan error)
	go func() {
		errors <- sess.Wait()
	}()
	fmt.Fprint(w, "\x00")
	w.Close()
	return b.Bytes(), false
}

func (tcloudcli *TcloudCli) UploadToUserDir(iscover bool, src string, dstDir string) (string, bool) {
	_, err := os.Stat(src)
	if err != nil {
		log.Printf("Failed to send to cluster. %s not exists.\n", src)
		return "", true
	}

	var cmd *exec.Cmd
	dst := tcloudcli.userConfig.SSHpath[0]
	dst = fmt.Sprintf("%s@%s:%s", tcloudcli.userConfig.UserName, dst, filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["userdir"], dstDir))
	ssh_config := fmt.Sprintf("/usr/bin/ssh -p %s -i %s", tcloudcli.userConfig.Port, tcloudcli.userConfig.AuthFile)
	if iscover {
		cmd = exec.Command("rsync", "-av", "--progress", "--delete", "--rsh", ssh_config, src, dst)
	} else {
		cmd = exec.Command("rsync", "-av", "--progress", "--rsh", ssh_config, src, dst)
	}

	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		log.Println("Failed to bind StdoutPipe")
		return dst, true
	}

	if err = cmd.Start(); err != nil {
		log.Println("Failed to start command")
		return dst, true
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}

	if err = cmd.Wait(); err != nil {
		return dst, true
	}
	return dst, false
}

func (tcloudcli *TcloudCli) UploadToWorkerDir(dirName string, src string) (string, bool) {
	_, err := os.Stat(src)
	if err != nil {
		log.Printf("Failed to send to cluster. %s not exists.\n", src)
		return "", true
	}

	dst := tcloudcli.userConfig.SSHpath[0]
	dst = fmt.Sprintf("%s@%s:%s", tcloudcli.userConfig.UserName, dst, filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["workdir"]))
	ssh_config := fmt.Sprintf("/usr/bin/ssh -p %s -i %s", tcloudcli.userConfig.Port, tcloudcli.userConfig.AuthFile)
	cmd := exec.Command("rsync", "-av", "--progress", "--delete", "--rsh", ssh_config, src, dst)

	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		return dst, true
	}

	if err = cmd.Start(); err != nil {
		return dst, true
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}

	if err = cmd.Wait(); err != nil {
		return dst, true
	}
	return dst, false
}

// SCP from SSHPath[0] to localhost
func (tcloudcli *TcloudCli) RecvFromCluster(src string, dst string, IsDir bool) bool {
	srcIP := tcloudcli.userConfig.SSHpath[0]
	srcPath := fmt.Sprintf("%s@%s:%s", tcloudcli.userConfig.UserName, srcIP, src)
	dstPath := fmt.Sprintf("%s", dst)

	var cmd *exec.Cmd
	if IsDir {
		cmd = exec.Command("scp", "-P", tcloudcli.userConfig.Port, "-r", "-i", tcloudcli.userConfig.AuthFile, srcPath, dstPath)
	} else {
		cmd = exec.Command("scp", "-P", tcloudcli.userConfig.Port, "-i", tcloudcli.userConfig.AuthFile, srcPath, dstPath)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Println("Failed to run cmd in RecvFromCluster ", err.Error(), stderr.String())
		return true
	}
	return false
}

func (tcloudcli *TcloudCli) BuildEnv(submitEnv *TACCGlobalEnv, args ...string) map[string]string {
	var config TuxivConfig
	localWorkDir, repoName, TACCDir, datasets, err := config.ParseTuxivConf(tcloudcli, submitEnv, args)
	randString := RandString(16)
	if err == true {
		log.Println("Failed to parse tuxiv.conf")
		os.Exit(-1)
	}
	var envName string
	if config.Environment.Name == "" {
		hashString := config.EnvNameGenerator()
		// fmt.Fprintln(w, fmt.Sprintf("name: %s", config.Environment.Name + "-" + hashString))
		envName = hashString
	} else {
		envName = config.Environment.Name
	}

	if err = tcloudcli.UploadRepo(repoName, localWorkDir); err == true {
		log.Println("Failed to upload env repository")
		os.Exit(-1)
	}

	if err = tcloudcli.AddSoftLink(datasets); err == true {
		log.Println("Failed to add softlink")
		os.Exit(-1)
	}

	// Remove auto-generated files
	if err = tcloudcli.RemoveAutoFiles(submitEnv); err == true {
		log.Println("Failed to remove all auto-generated files")
		os.Exit(-1)
	}

	// Generate env name and check if hit the cache, if so, return, otherwise, create new env.
	if tcloudcli.CondaCacheCheck(envName) {
		homeDir := filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName)
		condaBin := filepath.Join(homeDir, tcloudcli.clusterConfig.Conda)
		condaYaml := filepath.Join(homeDir, tcloudcli.clusterConfig.Dirs["workdir"], repoName, "configurations", "conda.yaml")
		cmd := fmt.Sprintf("%s %s env update -f %s -n %s\n", tcloudcli.prefix, condaBin, condaYaml, envName)
		if err := tcloudcli.RemoteExecCmd(cmd); err == true {
			log.Printf("Failed to update %s\n", envName)
			os.Exit(-1)
			return TACCDir
		}
		fmt.Printf("Env %s exists, dependencies updated.\n", envName)
		return TACCDir
	}
	if err = tcloudcli.CondaCreate(repoName, envName, randString); err == true {
		log.Println("Failed to create conda env")
		os.Exit(-1)
	}
	return TACCDir
}

func (tcloudcli *TcloudCli) UploadRepo(repoName string, localWorkDir string) bool {
	dst, err := tcloudcli.UploadToWorkerDir(repoName, localWorkDir)
	if err == true {
		log.Println("Failed to upload repo to ", dst)
		return true
	}
	// fmt.Println("Successfully upload repo to ", dst)
	return false
}

func (tcloudcli *TcloudCli) AddSoftLink(datasets []string) bool {
	for _, s := range datasets {
		cmd := fmt.Sprintf("curl -X GET %s/%s", CityNetAPI, s)

		out, err := tcloudcli.RemoteExecCmdOutput(cmd)
		if err == true {
			log.Println("Failed to access CityNet API")
			return true
		}

		var config Dataset
		json.Unmarshal(out, &config)

		datasetpath := filepath.Join(tcloudcli.clusterConfig.DatasetDir, config.Path)
		remoteUserDir := filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["userdir"])
		remoteDir := filepath.Join(remoteUserDir, config.Name)

		cmd = fmt.Sprintf("%s rm -f %s", tcloudcli.prefix, remoteDir)
		if err := tcloudcli.RemoteExecCmd(cmd); err == true {
			log.Printf("Failed to remove old softlink at %s\n", remoteDir)
			return true
		}
		cmd = fmt.Sprintf("%s ln -s %s %s", tcloudcli.prefix, datasetpath, remoteDir)
		if err := tcloudcli.RemoteExecCmd(cmd); err == true {
			log.Println("Failed to add softlink to user directory", err)
			return true
		}

		// fmt.Println("Successfully create softlink at", remoteDir)
	}
	return false
}

func (tcloudcli *TcloudCli) RemoveAutoFiles(submitEnv *TACCGlobalEnv) bool {
	if err := os.RemoveAll(submitEnv.LocalConfDir); err != nil {
		log.Println("Failed to remove auto-generated files")
		return true
	}
	if err := os.Remove(filepath.Join(submitEnv.LocalWorkDir, "run.sh")); err != nil {
		log.Println("Failed to remove auto-generated run.sh")
		return true
	}
	return false
}

func (tcloudcli *TcloudCli) CondaCreate(repoName string, envName string, randString string) bool {
	homeDir := filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName)
	condaBin := filepath.Join(homeDir, tcloudcli.clusterConfig.Conda)
	condaYaml := filepath.Join(homeDir, tcloudcli.clusterConfig.Dirs["workdir"], repoName, "configurations", "conda.yaml")
	cmd := fmt.Sprintf("%s %s env create -f %s -n %s\n", tcloudcli.prefix, condaBin, condaYaml, envName)
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		log.Println("Failed to run cmd in CondaCreate")
		return true
	}

	fmt.Printf("Successfully create environment %s\n.", envName)
	return false
}

func (tcloudcli *TcloudCli) CondaRemove(envName string) bool {
	homeDir := filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName)
	condaBin := filepath.Join(homeDir, tcloudcli.clusterConfig.Conda)
	cmd := fmt.Sprintf("%s %s remove -n %s --all -y", tcloudcli.prefix, condaBin, envName)
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		log.Println("Failed to run cmd in CondaRemove")
		return true
	}
	// fmt.Println("Previous environment \"", envName, "\" removed.")
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
	var submitEnv = NewGlobalEnv()

	cmd := fmt.Sprintf("mkdir -p  %s", filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["workdir"]))
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		log.Println("Failed to create remote workdir")
		return true
	}
	cmd = fmt.Sprintf("mkdir -p  %s", filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["userdir"]))
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		log.Println("Failed to create remote userdir")
		return true
	}
	cmd = fmt.Sprintf("mkdir -p  %s", filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["userdir"], submitEnv.SlurmUserlog))
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		log.Println("Failed to create remote workdir")
		return true
	}

	TACCDir := tcloudcli.BuildEnv(submitEnv, args...)

	cmd = fmt.Sprintf("%s sbatch %s", tcloudcli.prefix, filepath.Join(submitEnv.RemoteWorkDir, "configurations", "run.slurm"))

	// Create `RUNDIR` in remote and run cmd at `RUNDIR`
	cmd = fmt.Sprintf("mkdir -p %s && cd %s && %s", TACCDir["TACC_WORKDIR"], TACCDir["TACC_WORKDIR"], cmd)
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		log.Println("Failed to run cmd in tcloud submit")
		return true
	}
	fmt.Println("Job", submitEnv.RepoName, "submitted.")
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
		log.Println("Failed to run cmd in tcloud ps")
		return true
	}
	return false
}

func (tcloudcli *TcloudCli) XInit(args ...string) bool {
	// Remote receive config file
	src := DEFAULT_CLUSTERCONFIG_PATH
	dst := fmt.Sprintf("%s", filepath.Join(os.Getenv("HOME"), ".tcloud"))
	IsDir := false

	if err := tcloudcli.RecvFromCluster(src, dst, IsDir); err == true {
		log.Println("Failed to receive config files from TACC")
		return true
	}
	cmd := fmt.Sprintf("sinfo")
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		log.Println("Failed to get cluster information")
		return true
	}
	return false
}

func (tcloudcli *TcloudCli) XAdd(args ...string) bool {
	// Add new dependency to tuxiv.conf
	var config TuxivConfig
	err := config.AddDepTuxivFile(tcloudcli, args)
	if err == true {
		log.Println("Failed to add dependency to tuxiv.conf")
		os.Exit(-1)
	}
	return false
}

func (tcloudcli *TcloudCli) XInstall(args ...string) bool {
	var config TuxivConfig
	var submitEnv *TACCGlobalEnv
	_, _, _, _, err := config.ParseTuxivConf(tcloudcli, submitEnv, args)
	if err == true {
		log.Println("Failed to parse tuxiv.conf")
		os.Exit(-1)
	}
	condaYaml := filepath.Join(".", "configurations", "conda.yaml")
	removeCmd := exec.Command(tcloudcli.clusterConfig.Conda, "env", "remove", "-n", config.Environment.Name)
	if out, err := removeCmd.CombinedOutput(); err != nil {
		log.Println("Failed to create local environment. Err: ", err.Error())
		return true
	} else {
		fmt.Printf("%s\n", string(out))
	}
	createCmd := exec.Command(tcloudcli.clusterConfig.Conda, "env", "create", "-f", condaYaml)
	if out, err := createCmd.CombinedOutput(); err != nil {
		log.Println("Failed to create local environment. Err: ", err)
		return true
	} else {
		fmt.Printf("%s\n", string(out))
	}
	fmt.Println("Successfully create environment locally.")
	return false
}

func (tcloudcli *TcloudCli) XUpload(iscover bool, args ...string) bool {
	var src, dst string
	src = args[0]
	if len(args) > 1 {
		dst = args[1]
	} else {
		dst = "."
	}

	if _, err := tcloudcli.UploadToUserDir(iscover, src, dst); err {
		log.Printf("Failed to upload %s to %s\n", src, dst)
		return true
	}
	return false
}

func (tcloudcli *TcloudCli) XDownload(IsDir bool, args ...string) bool {
	var src, dst, remotesrc string
	src = args[0]
	if len(args) > 1 {
		dst = args[1]
	} else {
		dst = "."
	}

	// Format src, dst
	if src[0:1] == "./" {
		remotesrc = src[2:]
	} else if src[0] == '/' {
		remotesrc = src[1:]
	} else {
		remotesrc = src
	}

	remoteUserDir := filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["userdir"])
	remotesrc = filepath.Join(remoteUserDir, remotesrc)

	if err := tcloudcli.RecvFromCluster(remotesrc, dst, IsDir); err {
		if IsDir {
			log.Printf("Failed to receive directory %s to %s\n", src, dst)
			return true
		} else {
			log.Printf("Failed to receive file %s to %s\n", src, dst)
			return true
		}
	}
	return false
}

// Only allow remote workdir copy to remote userdir
// Src must contain repoName first
func (tcloudcli *TcloudCli) XCP(IsDir bool, args ...string) bool {
	var src, dst, remotesrc, remotedst string
	src = args[0]
	dst = args[1]

	// Format src, dst
	if src[0:1] == "./" {
		remotesrc = src[2:]
	} else if src[0] == '/' {
		remotesrc = src[1:]
	} else {
		remotesrc = src
	}

	if dst[0:1] == "./" {
		remotedst = dst[2:]
	} else if dst[0] == '/' {
		remotedst = dst[1:]
	} else {
		remotedst = dst
	}

	remoteWorkDir := filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["workdir"])
	remoteUserDir := filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["userdir"])
	remotesrc = filepath.Join(remoteWorkDir, remotesrc)
	remotedst = filepath.Join(remoteUserDir, remotedst)

	if IsDir {
		cmd := fmt.Sprintf("cp -r %s %s", remotesrc, remotedst)
		if err := tcloudcli.RemoteExecCmd(cmd); err == true {
			log.Printf("Failed to copy %s to %s\n", src, dst)
			return true
		}
	} else {
		cmd := fmt.Sprintf("mkdir -p %s && cp %s %s", remotedst, remotesrc, remotedst)
		if err := tcloudcli.RemoteExecCmd(cmd); err == true {
			log.Printf("Failed to copy %s to %s\n", src, dst)
			return true
		}
	}
	return false
}

func (tcloudcli *TcloudCli) XLS(IsLong bool, IsReverse bool, IsAll bool, args ...string) bool {
	var src, flags string
	if len(args) > 0 {
		src = args[0]
	} else {
		src = "."
	}
	flags = ""
	if IsLong {
		flags += " -l"
	}
	if IsReverse {
		flags += " -r"
	}
	if IsAll {
		flags += " -a"
	}

	remoteUserDir := filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["userdir"])
	remote := filepath.Join(remoteUserDir, src)

	cmd := fmt.Sprintf("ls %s %s", flags, remote)
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		log.Printf("Failed to ls %s %s\n", flags, remote)
		return true
	}
	return false
}

func (tcloudcli *TcloudCli) XENVLS(IsEnv bool, args ...string) bool {
	homeDir := filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName)
	condaBin := filepath.Join(homeDir, tcloudcli.clusterConfig.Conda)
	var src, flags string
	if len(args) > 0 {
		src = args[0]
	} else {
		src = ""
	}
	flags = ""
	if IsEnv {
		flags += "-n"
		cmd := fmt.Sprintf("%s list %s %s", condaBin, flags, src)
		if err := tcloudcli.RemoteExecCmd(cmd); err == true {
			log.Printf("Failed to list environment packages%s %s\n", flags, src)
			return true
		}
	} else {
		cmd := fmt.Sprintf("%s env list %s", condaBin, src)
		if err := tcloudcli.RemoteExecCmd(cmd); err == true {
			log.Printf("Failed to list environment %s %s\n", flags, src)
			return true
		}
	}
	return false
}

func (tcloudcli *TcloudCli) XCancel(job string, args ...string) bool {
	var cmd string
	cmd = fmt.Sprintf("%s scancel %s", tcloudcli.prefix, job)
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		log.Println("Failed to cancel job ", job)
		return true
	}
	return false
}

func (tcloudcli *TcloudCli) XDataset(args ...string) bool {
	if err := tcloudcli.AddSoftLink(args); err == true {
		log.Println("Failed to create dataset ", args[0])
		return true
	}
	return false
}

func (tcloudcli *TcloudCli) XCat(args ...string) bool {
	var cmd string
	remoteUserDir := filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["userdir"])
	remote := filepath.Join(remoteUserDir, args[0])

	cmd = fmt.Sprintf("cat %s", remote)
	if err := tcloudcli.RemoteExecCmd(cmd); err == true {
		log.Println("Failed to cat file ", args[0])
		return true
	}
	return false
}

func (tcloudcli *TcloudCli) CondaCacheCheck(envName string) bool {
	// Get env list from remote
	cmd := fmt.Sprintf("ls -ltr %s", filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, ".Miniconda3", "envs"))
	var envList []string
	if out, err := tcloudcli.RemoteExecCmdOutput(cmd); err == true {
		log.Println("Failed to get env list")
		os.Exit(1)
	} else {
		envList = strings.Split(strings.Trim(string(out), "\n"), "\n")
		for i, env := range envList {
			if i > 0 {
				splitString := strings.Split(env, " ")
				envList[i] = strings.Trim(splitString[len(splitString)-1], string(13))
			}
		}
		envList = envList[1:]
	}
	// Check if there is a hit, if so, return true, otherwise, return false
	for _, env := range envList {
		if env == envName {
			return true
		}
	}
	// Check the env cach length, if length > 10, remove the older env.
	envList = append(envList, envName)
	for {
		if len(envList) <= 10 {
			break
		}
		if err := tcloudcli.CondaRemove(envList[0]); err == true {
			log.Println("Failed to remove conda env")
			os.Exit(-1)
		}
		envList = envList[1:]
	}
	return false
}

func (tcloudcli *TcloudCli) XTest(args ...string) bool {
	var config TuxivConfig
	var src string
	src = args[0]

	err := config.DirSizeCheck(src, tcloudcli)
	log.Println("Error:", err)
	return false
}
