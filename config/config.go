package config

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/aaronriekenberg/go-httpd/logging"
)

var logger = logging.GetLogger()

type ResponseHeaders map[string]string

type BlockedLocation struct {
	ResponseStatus int `json:"responseStatus"`
}

type DirectoryLocation struct {
	StripPrefix   string `json:"stripPrefix"`
	DirectoryPath string `json:"directoryPath"`
}

type CompressedDirectoryLocation struct {
	StripPrefix   string `json:"stripPrefix"`
	DirectoryPath string `json:"directoryPath"`
}

type RedirectLocation struct {
	RedirectURL    string `json:"redirectURL"`
	ResponseStatus int    `json:"responseStatus"`
}

type FastCGILocation struct {
	Network        string `json:"network"`
	Address        string `json:"address"`
	ConnectionPool *struct {
		Size                 int `json:"size"`
		LifetimeMilliseconds int `json:"lifetimeMilliseconds"`
	} `json:"connectionPool"`
}

type Location struct {
	LocationID                  string                       `json:"locationID"`
	HttpPathPrefix              string                       `json:"httpPathPrefix"`
	BlockedLocation             *BlockedLocation             `json:"blockedLocation"`
	DirectoryLocation           *DirectoryLocation           `json:"directoryLocation"`
	CompressedDirectoryLocation *CompressedDirectoryLocation `json:"compressedDirectoryLocation"`
	RedirectLocation            *RedirectLocation            `json:"redirectLocation"`
	FastCGILocation             *FastCGILocation             `json:"fastCGILocation"`
	ResponseHeaders             *ResponseHeaders             `json:"responseHeaders"`
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

	logger.Printf("set httpServer.ReadTimeout = %v httpServer.WriteTimeout = %v", httpServer.ReadTimeout, httpServer.WriteTimeout)
}

type DropPrivileges struct {
	ChrootEnabled   bool   `json:"chrootEnabled"`
	ChrootDirectory string `json:"chrootDirectory"`
	GroupName       string `json:"groupName"`
	UserName        string `json:"userName"`
}

type RequestLogger struct {
	LogToStdout      bool   `json:"logToStdout"`
	RequestLogFile   string `json:"requestLogFile"`
	MaxSizeMegabytes int    `json:"maxSizeMegabytes"`
	MaxBackups       int    `json:"maxBackups"`
}

type NetworkAndListenAddress struct {
	Network       string `json:"network"`
	ListenAddress string `json:"listenAddress"`
}

type Server struct {
	ServerID                    string                    `json:"serverID"`
	NetworkAndListenAddressList []NetworkAndListenAddress `json:"networkAndListenAddressList"`
	TLSInfo                     *TLSInfo                  `json:"tlsInfo"`
	Timeouts                    *Timeouts                 `json:"timeouts"`
	ResponseHeaders             *ResponseHeaders          `json:"responseHeaders"`
	Locations                   []Location                `json:"locations"`
}

type Configuration struct {
	DropPrivileges *DropPrivileges `json:"dropPrivileges"`
	RequestLogger  *RequestLogger  `json:"requestLogger"`
	Servers        []Server        `json:"servers"`
}

func ReadConfiguration(configFile string) *Configuration {
	logger.Printf("reading json file %v", configFile)

	source, err := os.ReadFile(configFile)
	if err != nil {
		logger.Fatalf("error reading %v: %v", configFile, err)
	}

	var config Configuration
	if err = json.Unmarshal(source, &config); err != nil {
		logger.Fatalf("error parsing %v: %v", configFile, err)
	}

	return &config
}
