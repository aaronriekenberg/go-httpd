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

type serverRunFunc func(handler http.Handler)

var networkAndListenAddressToServerRunFunc = map[config.NetworkAndListenAddress]serverRunFunc{}

func createServer(
	serverConfig config.Server,
	networkAndListenAddress config.NetworkAndListenAddress,
) {

	if _, exists := networkAndListenAddressToServerRunFunc[networkAndListenAddress]; exists {
		logger.Fatalf("duplicate networkAndListenAddress %+v", networkAndListenAddress)
	}

	netListener, err := net.Listen(networkAndListenAddress.Network, networkAndListenAddress.ListenAddress)
	if err != nil {
		logger.Fatalf("net.Listen %+v: %v", networkAndListenAddress, err)
	}

	httpServer := &http.Server{
		Addr: networkAndListenAddress.ListenAddress,
	}

	serverConfig.Timeouts.ApplyToHTTPServer(httpServer)

	usingTLS := serverConfig.TLSInfo != nil

	if usingTLS {
		cert, err := tls.LoadX509KeyPair(serverConfig.TLSInfo.CertFile, serverConfig.TLSInfo.KeyFile)
		if err != nil {
			logger.Fatalf("Can't load certificates for server %v: %v", serverConfig.ServerID, err)
		}

		httpServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	networkAndListenAddressToServerRunFunc[networkAndListenAddress] = func(handler http.Handler) {

		httpServer.Handler = handler

		if usingTLS {

			logger.Printf("before ServeTLS serverID = %q networkAndListenAddress = %+v", serverConfig.ServerID, networkAndListenAddress)
			err := httpServer.ServeTLS(netListener, "", "")
			logger.Fatalf("server.ServeTLS err = %v serverID = %q networkAndListenAddress = %+v", err, serverConfig.ServerID, networkAndListenAddress)

		} else {

			logger.Printf("before Serve serverID = %q networkAndListenAddress = %+v", serverConfig.ServerID, networkAndListenAddress)
			err := httpServer.Serve(netListener)
			logger.Fatalf("server.Serve err = %v serverID = %q networkAndListenAddress = %+v", err, serverConfig.ServerID, networkAndListenAddress)

		}
	}

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

		handler := handlers.CreateServerLocationsHandler(
			serverConfig.Locations,
			serverConfig.CustomResponseHeaders,
		)

		handler = requestLogger.WrapHttpHandler(handler)

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

	runFunc, ok := networkAndListenAddressToServerRunFunc[networkAndListenAddress]
	if !ok {
		logger.Fatalf("unable to find runFunc networkAndListenAddress = %+v", networkAndListenAddress)
	}

	runFunc(handler)
}
