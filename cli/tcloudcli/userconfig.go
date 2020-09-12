package tcloudcli

// import (
// 	"encoding/json"
// 	"os"
// )

type UserConfig struct {
	// Slurm config
	authFile string
	userName string
}

func NewUserConfig(userName string, authFile string) *UserConfig {
	var config UserConfig
	config.authFile = authFile
	config.userName = userName
	return &config
}
