package servers

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/aaronriekenberg/go-httpd/handlers"
	"github.com/aaronriekenberg/go-httpd/logging"
	"github.com/aaronriekenberg/go-httpd/requestlogging"
)

var logger = logging.GetLogger()

type serverInfo struct {
	serverID    string
	netListener net.Listener
	httpServer  *http.Server
}

var networkAndListenAddressToServerInfo = map[config.NetworkAndListenAddress]*serverInfo{}

func createServer(
	serverConfig config.Server,
	networkAndListenAddress config.NetworkAndListenAddress,
) {
	if _, exists := networkAndListenAddressToServerInfo[networkAndListenAddress]; exists {
		logger.Fatalf("duplicate networkAndListenAddress %+v", networkAndListenAddress)
	}

	tcpListener, err := net.Listen(networkAndListenAddress.Network, networkAndListenAddress.ListenAddress)
	if err != nil {
		logger.Fatalf("net.Listen %+v: %v", networkAndListenAddress, err)
	}

	serverInfo := &serverInfo{
		serverID:    serverConfig.ServerID,
		netListener: tcpListener,
		httpServer: &http.Server{
			Addr: networkAndListenAddress.ListenAddress,
		},
	}

	serverConfig.Timeouts.ApplyToHTTPServer(serverInfo.httpServer)

	if serverConfig.TLSInfo != nil {
		cert, err := tls.LoadX509KeyPair(serverConfig.TLSInfo.CertFile, serverConfig.TLSInfo.KeyFile)
		if err != nil {
			logger.Fatalf("Can't load certificates for server %v: %v", serverConfig.ServerID, err)
		}

		serverInfo.httpServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	networkAndListenAddressToServerInfo[networkAndListenAddress] = serverInfo
}

func CreateServers(
	servers []config.Server,
) {
	logger.Printf("begin CreateServers")

	for _, serverConfig := range servers {
		for _, networkAndListenAddress := range serverConfig.NetworkAndListenAddressList {
			createServer(serverConfig, networkAndListenAddress)
		}
	}

	logger.Printf("end CreateServers")
}

func StartServers(
	servers []config.Server,
	requestLogger *requestlogging.RequestLogger,
) {
	logger.Printf("begin StartServers")

	for _, serverConfig := range servers {
		logger.Printf("StartServers serverID %q", serverConfig.ServerID)

		handler := handlers.CreateLocationsHandler(serverConfig.Locations)

		if requestLogger != nil {
			handler = requestLogger.WrapHttpHandler(handler)
		}

		for _, networkAndListenAddress := range serverConfig.NetworkAndListenAddressList {
			go runServer(
				networkAndListenAddress,
				handler,
			)
		}
	}

	logger.Printf("end StartServers")
}

func runServer(
	networkAndListenAddress config.NetworkAndListenAddress,
	handler http.Handler,
) {

	serverInfo, ok := networkAndListenAddressToServerInfo[networkAndListenAddress]
	if !ok {
		logger.Fatalf("unable to find serverInfo networkAndListenAddress = %+v", networkAndListenAddress)
	}

	httpServer := serverInfo.httpServer
	httpServer.Handler = handler

	if httpServer.TLSConfig != nil {

		logger.Printf("before ServeTLS serverID = %q networkAndListenAddress = %+v", serverInfo.serverID, networkAndListenAddress)
		err := httpServer.ServeTLS(serverInfo.netListener, "", "")
		logger.Fatalf("server.ServeTLS err = %v serverID = %q networkAndListenAddress = %+v", err, serverInfo.serverID, networkAndListenAddress)

	} else {

		logger.Printf("before Serve serverID = %q networkAndListenAddress = %+v", serverInfo.serverID, networkAndListenAddress)
		err := httpServer.Serve(serverInfo.netListener)
		logger.Fatalf("server.Serve err = %v serverID = %q networkAndListenAddress = %+v", err, serverInfo.serverID, networkAndListenAddress)

	}

}
