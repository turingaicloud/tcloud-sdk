package tcloudcli

import (
	"encoding/json"
	"os"
)

type UserConfig struct {
	UserName string   `json:"username"`
	SSHpath  []string `json:"sshpath"`
	AuthFile string   `json:"authfile"`
	Dir      []string `json:"dir"`
	path     string
}

func NewUserConfig(path string) *UserConfig {
	var config UserConfig
	file, err := os.Open(path)
	if err != nil {
		return &UserConfig{path: path}
	}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return &UserConfig{path: path}
	}

	config.path = path
	return &config
}
