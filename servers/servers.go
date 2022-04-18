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

var listenAddressToListeners = map[string]net.Listener{}

func CreateListeners(
	servers []config.Server,
) {
	log.Printf("begin CreateListeners")

	for _, serverConfig := range servers {
		for _, listenAddress := range serverConfig.ListenAddressList {

			if _, exists := listenAddressToListeners[listenAddress]; exists {
				log.Fatalf("duplicate listenAddress %q", listenAddress)
			}

			tcpListener, err := net.Listen("tcp", listenAddress)
			if err != nil {
				log.Fatalf("net.Listen %v: %v", listenAddress, err)
			}
			listenAddressToListeners[listenAddress] = tcpListener

			if serverConfig.TLSInfo != nil {
				cert, err := tls.LoadX509KeyPair(serverConfig.TLSInfo.CertFile, serverConfig.TLSInfo.KeyFile)
				if err != nil {
					log.Fatalf("Can't load certificates for server %v: %v", serverConfig.ServerID, err)
				}

				var tlsConfig tls.Config
				tlsConfig.Certificates = make([]tls.Certificate, 1)
				tlsConfig.Certificates[0] = cert

				tlsListener := tls.NewListener(tcpListener, &tlsConfig)

				listenAddressToListeners[listenAddress] = tlsListener
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

		for _, listenAddress := range serverConfig.ListenAddressList {
			go runServer(
				listenAddress,
				serverConfig,
				handler,
			)
		}
	}

	log.Printf("end StartServers")
}

func runServer(
	listenAddress string,
	serverConfig config.Server,
	handler http.Handler,
) {

	server := &http.Server{
		Addr:    listenAddress,
		Handler: handler,
	}

	serverConfig.Timeouts.ApplyToHTTPServer(server)

	listener, ok := listenAddressToListeners[listenAddress]
	if !ok {
		log.Fatalf("unable to find listener listenAddress = %v", listenAddress)
	}

	log.Printf("before Serve serverID = %q listenAddress = %q", serverConfig.ServerID, listenAddress)

	err := server.Serve(listener)

	log.Fatalf("server.Serve err = %v serverID = %q listenAddress = %q", err, serverConfig.ServerID, listenAddress)

}
