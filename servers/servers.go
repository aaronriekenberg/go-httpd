package servers

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"

	gorillaHandlers "github.com/gorilla/handlers"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/aaronriekenberg/go-httpd/handlers"
)

type serverInfo struct {
	listener  net.Listener
	tlsConfig *tls.Config
}

var networkAndListenAddressToServerInfo = map[config.NetworkAndListenAddress]*serverInfo{}

func CreateListeners(
	servers []config.Server,
) {
	log.Printf("begin CreateListeners")

	for _, serverConfig := range servers {
		for _, networkAndListenAddress := range serverConfig.NetworkAndListenAddressList {

			if _, exists := networkAndListenAddressToServerInfo[networkAndListenAddress]; exists {
				log.Fatalf("duplicate networkAndListenAddress %q", networkAndListenAddress)
			}

			serverInfo := &serverInfo{}

			tcpListener, err := net.Listen(networkAndListenAddress.Network, networkAndListenAddress.ListenAddress)
			if err != nil {
				log.Fatalf("net.Listen %+v: %v", networkAndListenAddress, err)
			}
			serverInfo.listener = tcpListener

			if serverConfig.TLSInfo != nil {
				cert, err := tls.LoadX509KeyPair(serverConfig.TLSInfo.CertFile, serverConfig.TLSInfo.KeyFile)
				if err != nil {
					log.Fatalf("Can't load certificates for server %v: %v", serverConfig.ServerID, err)
				}

				var tlsConfig tls.Config
				tlsConfig.Certificates = make([]tls.Certificate, 1)
				tlsConfig.Certificates[0] = cert

				tlsConfig.NextProtos = append(tlsConfig.NextProtos, "http/1.1", "h2")

				log.Printf("tlsConfig.NextProtos = %q", tlsConfig.NextProtos)

				tlsListener := tls.NewListener(tcpListener, &tlsConfig)

				serverInfo.listener = tlsListener
				serverInfo.tlsConfig = &tlsConfig
			}

			networkAndListenAddressToServerInfo[networkAndListenAddress] = serverInfo

		}
	}
}

func StartServers(
	servers []config.Server,
) {
	log.Printf("begin StartServers")

	for _, serverConfig := range servers {
		log.Printf("StartServers serverID %q", serverConfig.ServerID)

		handler := handlers.CreateLocationsHandler(serverConfig.Locations)

		if serverConfig.LogRequests {
			handler = gorillaHandlers.CombinedLoggingHandler(os.Stdout, handler)
		}

		for _, networkAndListenAddress := range serverConfig.NetworkAndListenAddressList {
			go runServer(
				networkAndListenAddress,
				serverConfig,
				handler,
			)
		}
	}

	log.Printf("end StartServers")
}

func runServer(
	networkAndListenAddress config.NetworkAndListenAddress,
	serverConfig config.Server,
	handler http.Handler,
) {

	serverInfo, ok := networkAndListenAddressToServerInfo[networkAndListenAddress]
	if !ok {
		log.Fatalf("unable to find serverInfo networkAndListenAddress = %+v", networkAndListenAddress)
	}

	server := &http.Server{
		Addr:      networkAndListenAddress.ListenAddress,
		Handler:   handler,
		TLSConfig: serverInfo.tlsConfig,
	}

	serverConfig.Timeouts.ApplyToHTTPServer(server)

	log.Printf("before Serve serverID = %q networkAndListenAddress = %+v", serverConfig.ServerID, networkAndListenAddress)

	err := server.Serve(serverInfo.listener)

	log.Fatalf("server.Serve err = %v serverID = %q networkAndListenAddress = %+v", err, serverConfig.ServerID, networkAndListenAddress)

}
