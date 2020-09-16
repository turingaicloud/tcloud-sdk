package tcloudcli

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
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

func (config *TuxivConfig) ParseTuxivConf(args []string) bool {
	var tuxivFile string
	var confDir string
	var workDir string
	if len(args) < 1 {
		tuxivFile = "tuxiv.conf"
		workDir = "."
		confDir = "./configurations"
	} else {
		tuxivFile = args[0]
		workDir = path.Dir(tuxivFile)
		confDir = filepath.Join(confDir, "configurations")
	}

	yamlFile, err := ioutil.ReadFile(tuxivFile)
	if err != nil {
		return true
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if _, err = os.Stat(confDir); os.IsNotExist(err) {
		os.Mkdir(confDir, 0755)
	}
	if err = config.CondaFile(workDir, confDir); err == true {
		fmt.Println("Environment config file generate failed.")
		return true
	}
	if err = config.SlurmFile(workDir, confDir); err == true {
		fmt.Println("Slurm config file generate failed.")
		return true
	}
	if err = config.CityFile(confDir); err == true {
		fmt.Println("Datasets config file generate failed.")
		return true
	}
	if err = config.RunshFile(workDir); err == true {
		fmt.Println("Run.sh exec file generate failed.")
		return true
	}
	return false
}

func (config *TuxivConfig) CondaFile(workDir string, confDir string) bool {
	f, err := os.Create(filepath.Join(confDir, "conda.yaml"))
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
	fmt.Fprintln(w, fmt.Sprintf("prefix: %s", filepath.Join(workDir, "environment")))
	w.Flush()
	return false
}

func (config *TuxivConfig) SlurmFile(workDir string, confDir string) bool {
	f, err := os.Create(filepath.Join(confDir, "run.slurm"))
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
	str := fmt.Sprintf("srun %s", filepath(workDir, "run.sh"))
	fmt.Fprintln(w, str)
	w.Flush()
	return false
}

func (config *TuxivConfig) CityFile(confDir string) bool {
	f, err := os.Create(filepath.Join(confDir, "citynet.sh"))
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

func (config *TuxivConfig) RunshFile(workDir string) bool {
	f, err := os.Create(filepath.Join(workDir, "run.sh"))
	if err != nil {
		fmt.Println("Create run.sh file failed.")
		return true
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	fmt.Fprintln(w, "#!/bin/bash\nsource ~/miniconda3/etc/profile.d/conda.sh")
	str := fmt.Sprintf("conda activate %s\n", config.Environment.Name)
	fmt.Fprintln(w, str)

	for _, s := range config.Entrypoint {
		str = fmt.Sprintf("%s \\", s)
		fmt.Fprintln(w, str)
	}
	w.Flush()
	if err = os.Chmod("go.sh", 0755); err != nil {
		fmt.Println("Run.sh chmod failed.")
		return true
	}
	return false
}
