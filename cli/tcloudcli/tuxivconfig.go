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
	}
	Datasets []struct {
		Name string
		Url  string
	}
}

func (config *TuxivConfig) TACCJobEnv(remoteWorkDir string) []string {
	var strlist []string
	// TACC Global Env
	strlist = append(strlist, fmt.Sprintf("TACC_WORKDIR=%s", remoteWorkDir))
	return strlist
}

func (config *TuxivConfig) ParseTuxivConf(tcloudcli *TcloudCli, args []string) (string, string, bool) {
	var tuxivFile string
	var localConfDir, localWorkDir string
	var remoteWorkDir string
	// var remoteConfDir string
	var repoName string
	if len(args) < 1 {
		tuxivFile = "tuxiv.conf"
		localWorkDir, _ = filepath.Abs(path.Dir(tuxivFile))
		localConfDir = filepath.Join(localWorkDir, "configurations")
		dirlist := strings.Split(localWorkDir, "/")
		repoName = dirlist[len(dirlist)-1]
		remoteWorkDir = fmt.Sprintf("/home/%s/%s", tcloudcli.userConfig.UserName, repoName)
		// remoteConfDir = filepath.Join(remoteWorkDir, "configurations")
	} else {
		tuxivFile = args[0]
		localWorkDir, _ = filepath.Abs(path.Dir(tuxivFile))
		localConfDir = filepath.Join(localWorkDir, "configurations")
		dirlist := strings.Split(localWorkDir, "/")
		repoName = dirlist[len(dirlist)-1]
		remoteWorkDir = fmt.Sprintf("/home/%s/%s", tcloudcli.userConfig.UserName, repoName)
		// remoteConfDir = filepath.Join(remoteWorkDir, "configurations")
	}

	yamlFile, err := ioutil.ReadFile(tuxivFile)
	if err != nil {
		return localWorkDir, repoName, true
	}

	err = yaml.Unmarshal(yamlFile, config)
	if _, err = os.Stat(localConfDir); os.IsNotExist(err) {
		os.Mkdir(localConfDir, 0755)
	}

	if err := config.CondaFile(localConfDir, remoteWorkDir); err == true {
		fmt.Println("Environment config file generate failed.")
		return localWorkDir, repoName, true
	}
	if err := config.SlurmFile(localConfDir, remoteWorkDir); err == true {
		fmt.Println("Slurm config file generate failed.")
		return localWorkDir, repoName, true
	}
	if err := config.CityFile(localConfDir); err == true {
		fmt.Println("Datasets config file generate failed.")
		return localWorkDir, repoName, true
	}
	if err := config.RunshFile(tcloudcli, localWorkDir); err == true {
		fmt.Println("Run.sh exec file generate failed.")
		return localWorkDir, repoName, true
	}
	return localWorkDir, repoName, false
}

func (config *TuxivConfig) CondaFile(localConfDir string, remoteWorkDir string) bool {
	f, err := os.Create(filepath.Join(localConfDir, "conda.yaml"))
	if err != nil {
		fmt.Println("Create Conda config file failed.")
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
	// prefix - set to ${workDir}/environment
	// fmt.Fprintln(w, fmt.Sprintf("prefix: %s", filepath.Join(remoteWorkDir, "environment")))
	w.Flush()
	return false
}

func (config *TuxivConfig) SlurmFile(localConfDir string, remoteWorkDir string) bool {
	f, err := os.Create(filepath.Join(localConfDir, "run.slurm"))
	if err != nil {
		fmt.Println("Create Slurm config file failed.")
		log.Fatal(err)
		return true
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	// Slurm file
	fmt.Fprintln(w, "#!/bin/bash")
	// SBATCH
	for _, s := range config.Job.General {
		str := fmt.Sprintf("#SBATCH --%s", s)
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
	for _, s := range config.TACCJobEnv(remoteWorkDir) {
		str := fmt.Sprintf("export %s", s)
		fmt.Fprintln(w, str)
	}
	str := fmt.Sprintf("srun %s", filepath.Join(remoteWorkDir, "run.sh"))
	fmt.Fprintln(w, str)
	w.Flush()
	return false
}

func (config *TuxivConfig) CityFile(localConfDir string) bool {
	f, err := os.Create(filepath.Join(localConfDir, "citynet.sh"))
	if err != nil {
		fmt.Println("Create Datasets config file failed.")
		return true
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	// TODO(Restructure cityFile) Not yet clear about the file format
	for _, s := range config.Datasets {
		str := fmt.Sprintf("%s\n%s\n", s.Name, s.Url)
		fmt.Fprintln(w, str)
	}
	w.Flush()
	return false
}

func (config *TuxivConfig) RunshFile(tcloudcli *TcloudCli, localWorkDir string) bool {
	f, err := os.Create(filepath.Join(localWorkDir, "run.sh"))
	if err != nil {
		fmt.Println("Create run.sh file failed.")
		return true
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	homeDir := fmt.Sprintf("/home/%s", tcloudcli.userConfig.UserName)
	str := fmt.Sprintf("#!/bin/bash\nsource %s/miniconda3/etc/profile.d/conda.sh", homeDir)
	fmt.Fprintln(w, str)
	str = fmt.Sprintf("conda activate %s\n", config.Environment.Name)
	fmt.Fprintln(w, str)

	for _, s := range config.Entrypoint {
		str = fmt.Sprintf("%s \\", s)
		fmt.Fprintln(w, str)
	}
	w.Flush()
	if err = os.Chmod("run.sh", 0755); err != nil {
		fmt.Println("Run.sh chmod failed.")
		return true
	}
	return false
}
func (config *TuxivConfig) AddDepTuxivFile(tcloudcli *TcloudCli, args []string) bool {
	var tuxivFile string
	tuxivFile = "tuxiv.conf"
	yamlFile, err := ioutil.ReadFile(tuxivFile)
	if err != nil {
		fmt.Println("Read file failed.")
		return true
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		fmt.Println("Parse original yaml file failed.")
		return true
	}
	for i := 0; i < len(config.Environment.Dependencies); i++ {
		slist := strings.Split(config.Environment.Dependencies[i], "=")
		deplist := strings.Split(args[0], "=")
		if deplist[0] == slist[0]{
			fmt.Println("Remove the original dependency %s", config.Environment.Dependencies[i])
			config.Environment.Dependencies = append(config.Environment.Dependencies[:i], config.Environment.Dependencies[i+1:]...)
		}
	}
	config.Environment.Dependencies = append(config.Environment.Dependencies, args[0])
	yamlFile, err = yaml.Marshal(config)
	if err != nil {
		fmt.Println("Format file failed.")
		return true
	}
	err = ioutil.WriteFile(tuxivFile, yamlFile, 0755)
	if err != nil {
		fmt.Println("Write file failed.")
		return true
	}

	return false
}
