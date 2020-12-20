package tcloudcli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type UserConfig struct {
	UserName string   `json: "username"`
	SSHpath  []string `json: "sshpath"`
	AuthFile string   `json: "authfile"`
	Port     string   `json: "port"`
	path     string   `json: "path"`
}

var DEFAULT_SSHPATH = "sing.cse.ust.hk"
var DEFAULT_AUTHFILE = ".ssh/id_rsa"
var DEFAULT_PORT = "30041"

func NewUserConfig(path string) *UserConfig {
	var config UserConfig
	file, err := os.Open(path)
	if err != nil {
		log.Println("Failed to open userconfig file")
		return &UserConfig{SSHpath: []string{DEFAULT_SSHPATH}, AuthFile: fmt.Sprintf("%s/%s", os.Getenv("HOME"), DEFAULT_AUTHFILE), Port: DEFAULT_PORT, path: path}
	}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return &UserConfig{SSHpath: []string{DEFAULT_SSHPATH}, AuthFile: fmt.Sprintf("%s/%s", os.Getenv("HOME"), DEFAULT_AUTHFILE), Port: DEFAULT_PORT, path: path}
	}

	// Set default value
	config.SSHpath = []string{DEFAULT_SSHPATH}
	if config.AuthFile == "" {
		config.AuthFile = fmt.Sprintf("%s/%s", os.Getenv("HOME"), DEFAULT_AUTHFILE)
	}
	config.Port = DEFAULT_PORT
	config.path = path
	return &config
}
