package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Entrypoint  []string
	Environment struct {
		Name       string
		Dependency []string
	}
	Job struct {
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
	},
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
	fmt.Println(setting)
	// Start parsing to 3 files
	conf_dir := "configurations"
	if _, err := os.Stat(conf_dir); os.IsNotExist(err) {
		os.Mkdir(conf_dir, 0755)
	}
	err1 := CondaFile(setting)
	err2 := SlurmFile(setting)
	err3 := CityFile(setting)
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

func CondaFile(setting Config) bool {
	// Switch map to output file

	return false
}

func SlurmFile(setting Config) bool {
	// Switch map to output file

	return false
}

func CityFile(setting Config) bool {
	// Switch map to output file

	return false
}
