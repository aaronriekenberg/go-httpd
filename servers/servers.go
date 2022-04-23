package servers

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/aaronriekenberg/go-httpd/handlers"
	"github.com/aaronriekenberg/go-httpd/requestlogger"
)

type serverInfo struct {
	serverID    string
	netListener net.Listener
	httpServer  *http.Server
}

var networkAndListenAddressToServerInfo = map[config.NetworkAndListenAddress]*serverInfo{}

func CreateServers(
	servers []config.Server,
) {
	log.Printf("begin CreateServers")

	for _, serverConfig := range servers {
		for _, networkAndListenAddress := range serverConfig.NetworkAndListenAddressList {

			if _, exists := networkAndListenAddressToServerInfo[networkAndListenAddress]; exists {
				log.Fatalf("duplicate networkAndListenAddress %q", networkAndListenAddress)
			}

			serverInfo := &serverInfo{
				serverID: serverConfig.ServerID,
				httpServer: &http.Server{
					Addr: networkAndListenAddress.ListenAddress,
				},
			}

			serverConfig.Timeouts.ApplyToHTTPServer(serverInfo.httpServer)

			tcpListener, err := net.Listen(networkAndListenAddress.Network, networkAndListenAddress.ListenAddress)
			if err != nil {
				log.Fatalf("net.Listen %+v: %v", networkAndListenAddress, err)
			}
			serverInfo.netListener = tcpListener

			if serverConfig.TLSInfo != nil {
				cert, err := tls.LoadX509KeyPair(serverConfig.TLSInfo.CertFile, serverConfig.TLSInfo.KeyFile)
				if err != nil {
					log.Fatalf("Can't load certificates for server %v: %v", serverConfig.ServerID, err)
				}

				tlsConfig := &tls.Config{
					Certificates: []tls.Certificate{cert},
					NextProtos:   []string{"h2", "http/1.1"},
				}

				tlsListener := tls.NewListener(tcpListener, tlsConfig)

				serverInfo.netListener = tlsListener
				serverInfo.httpServer.TLSConfig = tlsConfig
			}

			networkAndListenAddressToServerInfo[networkAndListenAddress] = serverInfo

		}
	}
}

func StartServers(
	servers []config.Server,
	requestLogger *requestlogger.RequestLogger,
) {
	log.Printf("begin StartServers")

	for _, serverConfig := range servers {
		log.Printf("StartServers serverID %q", serverConfig.ServerID)

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

	log.Printf("end StartServers")
}

func runServer(
	networkAndListenAddress config.NetworkAndListenAddress,
	handler http.Handler,
) {

	serverInfo, ok := networkAndListenAddressToServerInfo[networkAndListenAddress]
	if !ok {
		log.Fatalf("unable to find serverInfo networkAndListenAddress = %+v", networkAndListenAddress)
	}

	httpServer := serverInfo.httpServer
	httpServer.Handler = handler

	log.Printf("before Serve serverID = %q networkAndListenAddress = %+v", serverInfo.serverID, networkAndListenAddress)

	err := httpServer.Serve(serverInfo.netListener)

	log.Fatalf("server.Serve err = %v serverID = %q networkAndListenAddress = %+v", err, serverInfo.serverID, networkAndListenAddress)

}
