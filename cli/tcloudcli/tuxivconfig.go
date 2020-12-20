package tcloudcli

import (
	"bufio"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type TuxivConfig struct {
	Entrypoint  []string
	Environment struct {
		Name         string
		Channels     []string
		Dependencies []string
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
		submitEnv.RemoteWorkDir = fmt.Sprintf("%s/%s/%s/%s", tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["workdir"], submitEnv.RepoName)
		submitEnv.RemoteUserDir = fmt.Sprintf("%s/%s/%s", tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["userdir"])
	} else {
		submitEnv.LocalWorkDir, _ = filepath.Abs(path.Dir(args[0]))
		submitEnv.LocalConfDir = filepath.Join(submitEnv.LocalWorkDir, "configurations")
		dirlist := strings.Split(submitEnv.LocalWorkDir, "/")
		submitEnv.RepoName = dirlist[len(dirlist)-1]
		submitEnv.RemoteWorkDir = fmt.Sprintf("%s/%s/%s/%s", tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["workdir"], submitEnv.RepoName)
		submitEnv.RemoteUserDir = fmt.Sprintf("%s/%s/%s", tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName, tcloudcli.clusterConfig.Dirs["userdir"])
	}

	yamlFile, err := ioutil.ReadFile(tuxivFile)
	if err != nil {
		return submitEnv.LocalWorkDir, submitEnv.RepoName, TACCDir, nil, true
	}

	err = yaml.Unmarshal(yamlFile, config)
	if _, err = os.Stat(submitEnv.LocalConfDir); os.IsNotExist(err) {
		os.Mkdir(submitEnv.LocalConfDir, 0755)
	}

	if err := config.CondaFile(submitEnv.LocalConfDir, submitEnv.RemoteWorkDir); err == true {
		log.Println("Environment config file generate failed.")
		return submitEnv.LocalWorkDir, submitEnv.RepoName, TACCDir, nil, true
	}
	var err1 bool
	if TACCDir, err1 = config.SlurmFile(submitEnv, submitEnv.LocalConfDir, submitEnv.RemoteWorkDir, submitEnv.RemoteUserDir); err1 == true {
		log.Println("Slurm config file generate failed.")
		return submitEnv.LocalWorkDir, submitEnv.RepoName, TACCDir, config.Datasets, true
	}
	if err := config.CityFile(submitEnv.LocalConfDir); err == true {
		log.Println("Datasets config file generate failed.")
		return submitEnv.LocalWorkDir, submitEnv.RepoName, TACCDir, config.Datasets, true
	}
	if err := config.RunshFile(tcloudcli, submitEnv.LocalWorkDir); err == true {
		log.Println("Run.sh exec file generate failed.")
		return submitEnv.LocalWorkDir, submitEnv.RepoName, TACCDir, config.Datasets, true
	}
	return submitEnv.LocalWorkDir, submitEnv.RepoName, TACCDir, config.Datasets, false
}

func (config *TuxivConfig) CondaFile(localConfDir string, remoteWorkDir string) bool {
	f, err := os.Create(filepath.Join(localConfDir, "conda.yaml"))
	if err != nil {
		log.Println("Create Conda config file failed.")
		return true
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	// Conda file
	fmt.Fprintln(w, fmt.Sprintf("name: %s", config.Environment.Name))
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

func (config *TuxivConfig) SlurmFile(submitEnv *TACCGlobalEnv, localConfDir string, remoteWorkDir string, remoteUserDir string) (map[string]string, bool) {
	TACCDir := make(map[string]string)
	f, err := os.Create(filepath.Join(localConfDir, "run.slurm"))
	if err != nil {
		log.Println("Create Slurm config file failed.")
		log.Fatal(err)
		return TACCDir, true
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	// Slurm file
	fmt.Fprintln(w, "#!/bin/bash")
	// SBATCH
	for _, s := range config.Job.General {
		str := fmt.Sprintf("#SBATCH --%s", s)
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

func (config *TuxivConfig) CityFile(localConfDir string) bool {
	f, err := os.Create(filepath.Join(localConfDir, "citynet.sh"))
	if err != nil {
		log.Println("Create Datasets config file failed.")
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

func (config *TuxivConfig) RunshFile(tcloudcli *TcloudCli, localWorkDir string) bool {
	f, err := os.Create(filepath.Join(localWorkDir, "run.sh"))
	if err != nil {
		log.Println("Create run.sh file failed.")
		return true
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	homeDir := fmt.Sprintf("%s/%s", tcloudcli.clusterConfig.HomeDir, tcloudcli.userConfig.UserName)
	str := fmt.Sprintf("#!/bin/bash\nsource %s/%s", homeDir, CONDA_SHELL_PATH)
	fmt.Fprintln(w, str)
	str = fmt.Sprintf("conda activate %s\n", config.Environment.Name)
	fmt.Fprintln(w, str)

	for _, s := range config.Entrypoint {
		str = fmt.Sprintf("%s \\", s)
		fmt.Fprintln(w, str)
	}
	w.Flush()
	if err = os.Chmod("run.sh", 0755); err != nil {
		log.Println("Run.sh chmod failed.")
		return true
	}
	return false
}
func (config *TuxivConfig) AddDepTuxivFile(tcloudcli *TcloudCli, args []string) bool {
	var tuxivFile string
	tuxivFile = "tuxiv.conf"
	yamlFile, err := ioutil.ReadFile(tuxivFile)
	if err != nil {
		log.Println("Read file failed.")
		return true
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Println("Parse original yaml file failed.")
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
		log.Println("Format file failed.")
		return true
	}
	err = ioutil.WriteFile(tuxivFile, yamlFile, 0755)
	if err != nil {
		log.Println("Write file failed.")
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

	slurm_log := fmt.Sprintf("%s/%s", submitEnv.RemoteUserDir, submitEnv.SlurmUserlog)
	str = strings.Replace(str, "${TACC_SLURM_USERLOG}", slurm_log, -1)
	str = strings.Replace(str, "$TACC_SLURM_USERLOG", slurm_log, -1)
	return str
}
