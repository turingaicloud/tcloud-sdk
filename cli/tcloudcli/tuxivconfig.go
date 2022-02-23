package tcloudcli

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type TuxivConfig struct {
	Entrypoint  []string
	Environment struct {
		Name          string
		CachedEnvName string
		Channels      []string
		Dependencies  []string
	}
	Job struct {
		Name    string
		General []string
		Module  []string
		Env     []string
		Log     string
		// Model 	string
	}
	Datasets []string
}

var CONDA_SHELL_PATH = ".Miniconda3/etc/profile.d/conda.sh"

func (config *TuxivConfig) TACCJobEnv(remoteWorkDir string, remoteUserDir string) ([]string, map[string]string) {
	var strlist []string
	TACCDir := make(map[string]string)
	// TACC Global Env
	strlist = append(strlist, fmt.Sprintf("TACC_WORKDIR=%s", remoteWorkDir))
	TACCDir["TACC_WORKDIR"] = remoteWorkDir
	strlist = append(strlist, fmt.Sprintf("TACC_USERDIR=%s", remoteUserDir))
	TACCDir["TACC_USERDIR"] = remoteUserDir
	return strlist, TACCDir
}

func (config *TuxivConfig) ParseTuxivConf(tcloudcli *TcloudCli, submitEnv *TACCGlobalEnv, args []string) (string, string, map[string]string, []string, bool) {
	var tuxivFile = "tuxiv.conf"
	fmt.Println("Start parsing tuxiv.conf...")
	TACCDir := make(map[string]string)
	if len(args) < 1 {
		submitEnv.LocalWorkDir, _ = filepath.Abs(path.Dir("."))
		submitEnv.LocalConfDir = filepath.Join(submitEnv.LocalWorkDir, "configurations")
		dirlist := strings.Split(submitEnv.LocalWorkDir, "/")
		submitEnv.RepoName = dirlist[len(dirlist)-1]
		submitEnv.RemoteWorkDir = filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["workdir"], submitEnv.RepoName)
		submitEnv.RemoteUserDir = filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["userdir"])
	} else {
		submitEnv.LocalWorkDir, _ = filepath.Abs(path.Dir(args[0]))
		submitEnv.LocalConfDir = filepath.Join(submitEnv.LocalWorkDir, "configurations")
		dirlist := strings.Split(submitEnv.LocalWorkDir, "/")
		submitEnv.RepoName = dirlist[len(dirlist)-1]
		submitEnv.RemoteWorkDir = filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["workdir"], submitEnv.RepoName)
		submitEnv.RemoteUserDir = filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["userdir"])
	}
	yamlFile, err := ioutil.ReadFile(tuxivFile)
	if err != nil {
		return submitEnv.LocalWorkDir, submitEnv.RepoName, TACCDir, nil, true
	}

	err = yaml.Unmarshal(yamlFile, config)
	if _, err = os.Stat(submitEnv.LocalConfDir); os.IsNotExist(err) {
		os.Mkdir(submitEnv.LocalConfDir, 0755)
	}

	if err := config.CondaFile(submitEnv); err == true {
		log.Println("Failed to generate Environment config file")
		return submitEnv.LocalWorkDir, submitEnv.RepoName, TACCDir, nil, true
	}
	var err1 bool
	if TACCDir, err1 = config.SlurmFile(submitEnv); err1 == true {
		log.Println("Failed to generate Slurm config file")
		return submitEnv.LocalWorkDir, submitEnv.RepoName, TACCDir, config.Datasets, true
	}
	if err := config.CityFile(submitEnv); err == true {
		log.Println("Failed to generate Datasets config file")
		return submitEnv.LocalWorkDir, submitEnv.RepoName, TACCDir, config.Datasets, true
	}
	if err := config.RunshFile(tcloudcli, submitEnv); err == true {
		log.Println("Failed to generate Run.sh exec file")
		return submitEnv.LocalWorkDir, submitEnv.RepoName, TACCDir, config.Datasets, true
	}
	return submitEnv.LocalWorkDir, submitEnv.RepoName, TACCDir, config.Datasets, false
}

func (config *TuxivConfig) CondaFile(submitEnv *TACCGlobalEnv) bool {
	localConfDir := submitEnv.LocalConfDir
	f, err := os.Create(filepath.Join(localConfDir, "conda.yaml"))
	if err != nil {
		log.Println("Failed to create Conda config file")
		return true
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	// Conda file
	var EnvName string
	// TODO(wxc): verify if no cached_env_name is ""
	if config.Environment.CachedEnvName == "" {
		hashString := config.EnvNameGenerator()
		// fmt.Fprintln(w, fmt.Sprintf("name: %s", config.Environment.Name + "-" + hashString))
		EnvName = hashString
	} else {
		EnvName = config.Environment.Name
	}
	fmt.Fprintln(w, fmt.Sprintf("name: %s", EnvName))
	// Channels
	fmt.Fprintln(w, fmt.Sprintf("channels:"))
	for _, s := range config.Environment.Channels {
		str := fmt.Sprintf("  - %s", s)
		fmt.Fprintln(w, str)
	}
	// Dependencies
	fmt.Fprintln(w, "dependencies:")
	for _, s := range config.Environment.Dependencies {
		str := fmt.Sprintf("  - %s", s)
		fmt.Fprintln(w, str)
	}
	w.Flush()
	return false
}

func (config *TuxivConfig) SlurmFile(submitEnv *TACCGlobalEnv) (map[string]string, bool) {
	localConfDir := submitEnv.LocalConfDir
	remoteWorkDir := submitEnv.RemoteWorkDir
	remoteUserDir := submitEnv.RemoteUserDir
	TACCDir := make(map[string]string)
	f, err := os.Create(filepath.Join(localConfDir, "run.slurm"))
	if err != nil {
		log.Println("Failed to create Slurm config file")
		log.Fatal(err)
		return TACCDir, true
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	// Slurm file
	fmt.Fprintln(w, "#!/bin/bash")
	// SBATCH
	var CHECK_OUTPUT = false
	for _, s := range config.Job.General {
		if strings.Contains(s, "output=") == true {
			CHECK_OUTPUT = true
		}
		str := fmt.Sprintf("#SBATCH --%s", s)
		str = ReplaceGlobalEnv(str, submitEnv)
		fmt.Fprintln(w, str)
	}
	if CHECK_OUTPUT == false {
		str := fmt.Sprintf("#SBATCH --output=%s", "${TACC_SLURM_USERLOG}/slurm-%j.out")
		str = ReplaceGlobalEnv(str, submitEnv)
		fmt.Fprintln(w, str)
	}
	// Module
	for _, s := range config.Job.Module {
		str := fmt.Sprintf("module load %s", s)
		fmt.Fprintln(w, str)
	}
	// Env
	for _, s := range config.Job.Env {
		str := fmt.Sprintf("export %s", s)
		fmt.Fprintln(w, str)
	}

	// TACC Env
	strlist, TACCDir := config.TACCJobEnv(remoteWorkDir, remoteUserDir)
	for _, s := range strlist {
		str := fmt.Sprintf("export %s", s)
		fmt.Fprintln(w, str)
	}
	str := fmt.Sprintf("srun %s", filepath.Join(remoteWorkDir, "run.sh"))
	fmt.Fprintln(w, str)
	w.Flush()
	return TACCDir, false
}

func (config *TuxivConfig) CityFile(submitEnv *TACCGlobalEnv) bool {
	localConfDir := submitEnv.LocalConfDir
	f, err := os.Create(filepath.Join(localConfDir, "citynet.sh"))
	if err != nil {
		log.Println("Failed to create Datasets config file")
		return true
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, s := range config.Datasets {
		fmt.Fprintln(w, s)
	}

	w.Flush()
	return false
}

func (config *TuxivConfig) RunshFile(tcloudcli *TcloudCli, submitEnv *TACCGlobalEnv) bool {
	localWorkDir := submitEnv.LocalWorkDir
	f, err := os.Create(filepath.Join(localWorkDir, "run.sh"))
	if err != nil {
		log.Println("Failed to create run.sh file")
		return true
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	homeDir := filepath.Join(tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName)
	str := fmt.Sprintf("#!/bin/bash\nsource %s/%s", homeDir, CONDA_SHELL_PATH)
	fmt.Fprintln(w, str)

	var EnvName string
	// TODO(wxc): verify if no cached_env_name is ""
	if config.Environment.CachedEnvName == "" {
		hashString := config.EnvNameGenerator()
		// fmt.Fprintln(w, fmt.Sprintf("name: %s", config.Environment.Name + "-" + hashString))
		EnvName = hashString
	} else {
		EnvName = config.Environment.Name
	}
	str = fmt.Sprintf("conda activate %s\n", EnvName)
	fmt.Fprintln(w, str)

	for _, s := range config.Entrypoint {
		str = fmt.Sprintf("%s \\", s)
		fmt.Fprintln(w, str)
	}
	w.Flush()
	if err = os.Chmod("run.sh", 0755); err != nil {
		log.Println("Failed to chmod run.sh")
		return true
	}
	return false
}
func (config *TuxivConfig) AddDepTuxivFile(tcloudcli *TcloudCli, args []string) bool {
	var tuxivFile string
	tuxivFile = "tuxiv.conf"
	yamlFile, err := ioutil.ReadFile(tuxivFile)
	if err != nil {
		log.Println("Failed to read file")
		return true
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Println("Failed to parse original yaml file")
		return true
	}
	for i := 0; i < len(config.Environment.Dependencies); i++ {
		slist := strings.Split(config.Environment.Dependencies[i], "=")
		deplist := strings.Split(args[0], "=")
		if deplist[0] == slist[0] {
			fmt.Println("Remove the original dependency", config.Environment.Dependencies[i])
			config.Environment.Dependencies = append(config.Environment.Dependencies[:i], config.Environment.Dependencies[i+1:]...)
		}
	}
	config.Environment.Dependencies = append(config.Environment.Dependencies, args[0])
	yamlFile, err = yaml.Marshal(config)
	if err != nil {
		log.Println("Failed to format file")
		return true
	}
	err = ioutil.WriteFile(tuxivFile, yamlFile, 0755)
	if err != nil {
		log.Println("Failed to write file")
		return true
	}

	return false
}

// Currently only replace TACC_WORKDIR, TACC_USERDIR
func ReplaceGlobalEnv(str string, submitEnv *TACCGlobalEnv) string {
	str = strings.Replace(str, "${TACC_WORKDIR}", submitEnv.RemoteWorkDir, -1)
	str = strings.Replace(str, "$TACC_WORKDIR", submitEnv.RemoteWorkDir, -1)
	str = strings.Replace(str, "${TACC_USERDIR}", submitEnv.RemoteUserDir, -1)
	str = strings.Replace(str, "$TACC_USERDIR", submitEnv.RemoteUserDir, -1)

	slurm_log := filepath.Join(submitEnv.RemoteUserDir, submitEnv.SlurmUserlog)
	str = strings.Replace(str, "${TACC_SLURM_USERLOG}", slurm_log, -1)
	str = strings.Replace(str, "$TACC_SLURM_USERLOG", slurm_log, -1)
	return str
}
func (config *TuxivConfig) EnvNameGenerator() string {
	// Parse package (with version) list from conda.yaml
	dep := config.Environment.Dependencies
	// Sort the package by Alphabetical order and contact as a string
	sort.Strings(dep)
	// Generate md5 hash value from the string
	jointDep := strings.Join(dep, " ")
	data := []byte(jointDep)
	hashValue := md5.Sum(data)
	hashString := hex.EncodeToString(hashValue[:])
	return hashString
}
