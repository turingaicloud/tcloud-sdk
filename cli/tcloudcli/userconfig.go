package tcloudcli

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type UserConfig struct {
	UserName string   `json: "username"`
	SSHpath  []string `json: "sshpath"`
	AuthFile string   `json: "authfile"`
	Port     string   `json: "port"`
	path     string   `json: "path"`
}

var DEFAULT_SSHPATH = "gw.turing.ust.hk"
var DEFAULT_AUTHFILE = ".ssh/id_rsa"
var DEFAULT_PORT = "30041"

func NewUserConfig(path string) *UserConfig {
	var config UserConfig
	file, err := os.Open(path)
	if err != nil {
		log.Println("Failed to open userconfig file")
		return &UserConfig{SSHpath: []string{DEFAULT_SSHPATH}, AuthFile: filepath.Join(os.Getenv("HOME"), DEFAULT_AUTHFILE), Port: DEFAULT_PORT, path: path}
	}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		// decode local config file failed, using default value instead
		return &UserConfig{SSHpath: []string{DEFAULT_SSHPATH}, AuthFile: filepath.Join(os.Getenv("HOME"), DEFAULT_AUTHFILE), Port: DEFAULT_PORT, path: path}
	}
	return &config
}
