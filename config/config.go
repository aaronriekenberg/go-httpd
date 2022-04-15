package config

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type BlockedLocation struct {
	ResponseStatus int `json:"responseStatus"`
}

type DirectoryLocation struct {
	StripPrefix       string `json:"stripPrefix"`
	CacheControlValue string `json:"cacheControlValue"`
	DirectoryPath     string `json:"directoryPath"`
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

type Timeouts struct {
	ReadTimeoutMilliseconds  int `json:"readTimeoutMilliseconds"`
	WriteTimeoutMilliseconds int `json:"writeTimeoutMilliseconds"`
}

func (timeouts *Timeouts) ApplyToHTTPServer(httpServer *http.Server) {
	if timeouts == nil {
		return
	}

	httpServer.ReadTimeout = time.Duration(timeouts.ReadTimeoutMilliseconds) * time.Millisecond
	httpServer.WriteTimeout = time.Duration(timeouts.WriteTimeoutMilliseconds) * time.Millisecond

	log.Printf("set httpServer.ReadTimeout = %v httpServer.WriteTimeout = %v", httpServer.ReadTimeout, httpServer.WriteTimeout)
}

type Server struct {
	ServerID          string     `json:"serverID"`
	ListenAddressList []string   `json:"listenAddressList"`
	TLSInfo           *TLSInfo   `json:"tlsInfo"`
	Timeouts          *Timeouts  `json:"timeouts"`
	LogRequests       bool       `json:"logRequests"`
	Locations         []Location `json:"locations"`
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
