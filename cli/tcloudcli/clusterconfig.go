package tcloudcli

import (
	"encoding/json"
	"log"
	"os"
)

type ClusterConfig struct {
	Dirs          map[string]string `json: "dirs"`
	TcloudVersion string            `json: "tcloudversion"`
	HomeDir       string            `json: "homedir"`
	DatasetDir    string            `json: "datasetdir"`
	Conda         string            `json: "conda"`
	StorageQuota  int64             `json: "storagequota"`
	// Note: StorageQuota number in clusterconfig is in MB
	path string `json: "path"`
}

func NewClusterConfig(path string) *ClusterConfig {
	var config ClusterConfig
	file, err := os.Open(path)
	if err != nil {
		log.Println("Failed to open clusterconfig file")
		return &ClusterConfig{path: path}
	}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return &ClusterConfig{path: path}
	}
	config.path = path
	return &config
}
