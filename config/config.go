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

type RedirectLocation struct {
	RedirectURL    string `json:"redirectURL"`
	ResponseStatus int    `json:"responseStatus"`
}

type FastCGILocation struct {
	UnixSocketPath string `json:"unixSocketPath"`
}

type Location struct {
	HttpPathPrefix    string             `json:"httpPathPrefix"`
	BlockedLocation   *BlockedLocation   `json:"blockedLocation"`
	DirectoryLocation *DirectoryLocation `json:"directoryLocation"`
	RedirectLocation  *RedirectLocation  `json:"redirectLocation"`
	FastCGILocation   *FastCGILocation   `json:"fastCGILocation"`
}

type TLSInfo struct {
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

type Server struct {
	ListenAddress string     `json:"listenAddress"`
	Locations     []Location `json:"locations"`
	TLSInfo       *TLSInfo   `json:"tlsInfo"`
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
