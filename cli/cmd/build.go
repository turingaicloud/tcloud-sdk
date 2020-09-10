package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
        "os/exec"
)

type Config struct {
	Entrypoint  []string
	Environment struct {
		Name         string
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

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Setup conda environment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tcloud build CLI")
		err := ParseTuxivConf(args)
		if err {
			fmt.Println("Parse tuxiv config file failed.")
			log.Fatal(err)
		}
                err_ := CondaCreate()
                if err_ {
                        fmt.Println("Create conda failed")
                        log.Fatal(err_)
		}
	},
}
func CondaCreate() bool {
    cmd := exec.Command("conda", "env", "create", "-f", "configurations/conda.yaml")
    out, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
        return true
    }
    fmt.Printf("combined out:\n%s\n", string(out))
    return false
}
func ParseTuxivConf(args []string) bool {
	// Check tuxiv.conf
	tuxivFile := "tuxiv.conf"
	if len(args) > 0 {
		tuxivFile = args[0]
	}

	yamlFile, err := ioutil.ReadFile(tuxivFile)
	if err != nil {
		log.Fatal(err)
	}
	var setting Config
	// resultMap := make(map[string]interface{})
	err = yaml.Unmarshal(yamlFile, &setting)
	// Start parsing to 3 files
	conf_dir := "configurations"
	if _, err := os.Stat(conf_dir); os.IsNotExist(err) {
		os.Mkdir(conf_dir, 0755)
	}
	err1 := CondaFile(setting, conf_dir)
	err2 := SlurmFile(setting, conf_dir)
	err3 := CityFile(setting, conf_dir)
	if err1 || err2 || err3 {
		if err1 {
			fmt.Println("Environment config file generate failed.")
			log.Fatal(err1)
		}
		if err2 {
			fmt.Println("Slurm config file generate failed.")
			log.Fatal(err2)
		}
		if err2 {
			fmt.Println("Datasets config file generate failed.")
			log.Fatal(err1)
		}
	}
	return err1 || err2 || err3
}

func CondaFile(setting Config, conf_dir string) bool {
	f, err := os.Create(conf_dir + "/conda.yaml")
	if err != nil {
		fmt.Println("Create Conda config file failed.")
		log.Fatal(err)
		return true
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	// Conda file
	fmt.Fprintln(w, fmt.Sprintf("name: %s", setting.Environment.Name))
        // Channels
        fmt.Fprintln(w, fmt.Sprintf("channels:\n  - defaults"))
	// Dependencies
	fmt.Fprintln(w, "dependencies:")
	for _, s := range setting.Environment.Dependencies {
		str := fmt.Sprintf("  - %s", s)
		fmt.Fprintln(w, str)
	}
        // Prefix
        fmt.Fprintln(w, fmt.Sprint("prefix: ../environment"))
	w.Flush()
	return false
}

func SlurmFile(setting Config, conf_dir string) bool {
	f, err := os.Create(conf_dir + "/run.slurm")
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
	for _, s := range setting.Job.General {
		str := fmt.Sprintf("#SBATCH --%s", s)
		fmt.Fprintln(w, str)
	}
	// Module
	for _, s := range setting.Job.Module {
		str := fmt.Sprintf("module load %s", s)
		fmt.Fprintln(w, str)
	}
	// Env
	for _, s := range setting.Job.Env {
		str := fmt.Sprintf("export %s", s)
		fmt.Fprintln(w, str)
	}
	w.Flush()
	return false
}

func CityFile(setting Config, conf_dir string) bool {
	f, err := os.Create(conf_dir + "/citynet.sh")
	if err != nil {
		fmt.Println("Create Datasets config file failed.")
		log.Fatal(err)
		return true
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	// CityNet file
	// Not yet clear about the file format
	for _, s := range setting.Datasets {
		str := fmt.Sprintf("\n%s\n%s\n", s.Name, s.Url)
		fmt.Fprintln(w, str)
	}
	w.Flush()
	return false
}
