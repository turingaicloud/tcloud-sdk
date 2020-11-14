package tcloudcli

import (
	"encoding/json"
	"fmt"
	"os"
)

type UserConfig struct {
	UserName string   `json: "username"`
	SSHpath  []string `json: "sshpath"`
	AuthFile string   `json: "authfile"`
	path     string   `json: "path"`
}

func NewUserConfig(path string) *UserConfig {
	var config UserConfig
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Failed to open userconfig file")
		return &UserConfig{path: path}
	}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		fmt.Println("Failed to parse userconfig file")
		return &UserConfig{path: path}
	}

	config.path = path
	return &config
}
