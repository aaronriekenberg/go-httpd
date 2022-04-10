package config

import (
	"encoding/json"
	"log"
	"os"
)

type BlockedLocation struct {
	ResponseStatus int `json:"responseStatus"`
}

type DirectoryLocation struct {
	StripPrefix   string `json:"stripPrefix"`
	DirectoryPath string `json:"directoryPath"`
}

type Location struct {
	HttpPathPrefix    string             `json:"httpPathPrefix"`
	BlockedLocation   *BlockedLocation   `json:"blockedLocation"`
	DirectoryLocation *DirectoryLocation `json:"directoryLocation"`
}

type Server struct {
	ListenAddress string     `json:"listenAddress"`
	Locations     []Location `json:"locations"`
}

type Configuration struct {
	Servers []Server `json:"servers"`
}

func ReadConfiguration(configFile string) *Configuration {
	log.Printf("reading json file %v", configFile)

	source, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error reading %v: %v", configFile, err)
	}

	var config Configuration
	if err = json.Unmarshal(source, &config); err != nil {
		log.Fatalf("error parsing %v: %v", configFile, err)
	}

	return &config
}
