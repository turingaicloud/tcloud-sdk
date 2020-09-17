package tcloudcli

import (
	"encoding/json"
	"os"
)

type UserConfig struct {
	UserName string   `json:"username"`
	SSHpath  []string `json:"ssh"`
	path     string
	authFile string
}

func NewUserConfig(path string, authFile string) *UserConfig {
	var config UserConfig
	file, err := os.Open(path)
	if err != nil {
		return &UserConfig{path: path, authFile: authFile}
	}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return &UserConfig{path: path, authFile: authFile}
	}

	config.path = path
	config.authFile = authFile
	return &config
}
