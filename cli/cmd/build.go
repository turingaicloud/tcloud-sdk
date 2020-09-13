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
	"tcloud-sdk/cli/tcloudcli"
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

func NewBuildCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Parse tuxiv.confg and Setup conda environment",
		Run: func(cmd *cobra.Command, args []string) {
			// fmt.Println("tcloud build CLI")
			setting, err := ParseTuxivConf(args)
			if err {
				fmt.Println("Parse tuxiv config file failed.")
				log.Fatal(err)
			}
			err_0 := UploadFile()
			if err_0 {
				fmt.Println("Upload file env failed")
				log.Fatal(err_0)
			}
			err_1 := CondaRemove(setting.Environment.Name)
			if err_1 {
				fmt.Println("remove conda env failed")
				log.Fatal(err_1)
			}
			err_2 := CondaCreate()
			if err_2 {
				fmt.Println("Create conda env failed")
				log.Fatal(err_2)
			}
		},
	}
}
func UploadFile() bool {
	cmd := exec.Command("scp", "-r", "-i", "./TACC.pem", "../tcloud_job", "ubuntu@18.162.45.250:/home/ubuntu")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Printf("Upload file to TACC JUMP out:\n%s\n", string(out))
	bash_command := "scp -r /home/ubuntu/tcloud_job TACC1:/home/ubuntu"
	cmd = exec.Command("ssh", "-i", "./TACC.pem", "ubuntu@18.162.45.250", bash_command)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Printf("Upload file to TACC1 out:\n%s\n", string(out))
	bash_command = "scp -r /home/ubuntu/tcloud_job TACC2:/home/ubuntu"
	cmd = exec.Command("ssh", "-i", "./TACC.pem", "ubuntu@18.162.45.250", bash_command)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Printf("Upload file to TACC2 out:\n%s\n", string(out))
	return false
}
func CondaCreate() bool {
	bash_command := "ssh TACC1 /home/ubuntu/miniconda3/bin/conda env create -f /home/ubuntu/tcloud_job/configurations/conda.yaml"
	cmd := exec.Command("ssh", "-i", "./TACC.pem", "ubuntu@18.162.45.250", bash_command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Printf("conda create on TACC1 out:\n%s\n", string(out))
	bash_command = "ssh TACC2 /home/ubuntu/miniconda3/bin/conda env create -f /home/ubuntu/tcloud_job/configurations/conda.yaml"
	cmd = exec.Command("ssh", "-i", "./TACC.pem", "ubuntu@18.162.45.250", bash_command)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Printf("conda create on TACC2 out:\n%s\n", string(out))
	return false
}
func CondaRemove(name string) bool {
	bash_command := "ssh TACC1 /home/ubuntu/miniconda3/bin/conda remove -n " + name + " --all -y"
	cmd := exec.Command("ssh", "-i", "./TACC.pem", "ubuntu@18.162.45.250", bash_command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Printf("conda remove on TACC1 out:\n%s\n", string(out))
	bash_command = "ssh TACC2 /home/ubuntu/miniconda3/bin/conda remove -n " + name + " --all -y"
	cmd = exec.Command("ssh", "-i", "./TACC.pem", "ubuntu@18.162.45.250", bash_command)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	fmt.Printf("conda remove on TACC2 out:\n%s\n", string(out))
	return false
}
func ParseTuxivConf(args []string) (Config, bool) {
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
	err4 := GoshFile(setting)
	err5 := os.Chmod("go.sh", 0755)
	if err5 != nil {
		log.Println(err5)
	}
	if err1 || err2 || err3 || err4 {
		if err1 {
			fmt.Println("Environment config file generate failed.")
			log.Fatal(err1)
		}
		if err2 {
			fmt.Println("Slurm config file generate failed.")
			log.Fatal(err2)
		}
		if err3 {
			fmt.Println("Datasets config file generate failed.")
			log.Fatal(err3)
		}
		if err4 {
			fmt.Println("Datasets config file generate failed.")
			log.Fatal(err4)
		}
	}
	return setting, err1 || err2 || err3 || err4
}
func GoshFile(setting Config) bool {
	f, err := os.Create("./go.sh")
	if err != nil {
		fmt.Println("Create go.sh file failed.")
		log.Fatal(err)
		return true
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	fmt.Fprintln(w, "#!/bin/bash\nsource ~/miniconda3/etc/profile.d/conda.sh\nconda activate tf\n")
	for _, s := range setting.Entrypoint {
		str := fmt.Sprintf("%s \\", s)
		fmt.Fprintln(w, str)
	}
	w.Flush()
	return false
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
	str := fmt.Sprintf("srun -N2 ../go.sh")
	fmt.Fprintln(w, str)
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
