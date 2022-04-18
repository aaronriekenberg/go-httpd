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

var networkAndListenAddressToListeners = map[config.NetworkAndListenAddress]net.Listener{}

func CreateListeners(
	servers []config.Server,
) {
	log.Printf("begin CreateListeners")

	for _, serverConfig := range servers {
		for _, networkAndListenAddress := range serverConfig.NetworkAndListenAddressList {

			if _, exists := networkAndListenAddressToListeners[networkAndListenAddress]; exists {
				log.Fatalf("duplicate networkAndListenAddress %q", networkAndListenAddress)
			}

			tcpListener, err := net.Listen(networkAndListenAddress.Network, networkAndListenAddress.ListenAddress)
			if err != nil {
				log.Fatalf("net.Listen %+v: %v", networkAndListenAddress, err)
			}
			networkAndListenAddressToListeners[networkAndListenAddress] = tcpListener

			if serverConfig.TLSInfo != nil {
				cert, err := tls.LoadX509KeyPair(serverConfig.TLSInfo.CertFile, serverConfig.TLSInfo.KeyFile)
				if err != nil {
					log.Fatalf("Can't load certificates for server %v: %v", serverConfig.ServerID, err)
				}

				var tlsConfig tls.Config
				tlsConfig.Certificates = make([]tls.Certificate, 1)
				tlsConfig.Certificates[0] = cert

				tlsListener := tls.NewListener(tcpListener, &tlsConfig)

				networkAndListenAddressToListeners[networkAndListenAddress] = tlsListener
			}

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

	server := &http.Server{
		Addr:    networkAndListenAddress.ListenAddress,
		Handler: handler,
	}

	serverConfig.Timeouts.ApplyToHTTPServer(server)

	listener, ok := networkAndListenAddressToListeners[networkAndListenAddress]
	if !ok {
		log.Fatalf("unable to find listener networkAndListenAddress = %+v", networkAndListenAddress)
	}

	log.Printf("before Serve serverID = %q networkAndListenAddress = %+v", serverConfig.ServerID, networkAndListenAddress)

	err := server.Serve(listener)

	log.Fatalf("server.Serve err = %v serverID = %q networkAndListenAddress = %+v", err, serverConfig.ServerID, networkAndListenAddress)

}
