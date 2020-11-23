package tcloudcli

import (
	"encoding/json"
	"fmt"
	"os"
)

type ClusterConfig struct {
	Dirs          map[string]string `json:"dirs"`
	TcloudVersion string            `json:"tcloudversion"`
	path          string            `json: "path"`
}

func NewClusterConfig(path string) *ClusterConfig {
	var config ClusterConfig
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Failed to open clusterconfig file")
		return &ClusterConfig{path: path}
	}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		// fmt.Println("Failed to parse clusterconfig file")
		return &ClusterConfig{path: path}
	}
	// fmt.Println(config)
	config.path = path
	return &config
}
