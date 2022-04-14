package servers

import (
	"log"
	"net/http"
	"os"

	gorillaHandlers "github.com/gorilla/handlers"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/aaronriekenberg/go-httpd/handlers"
)

func runServer(
	listenAddress string,
	serverConfig config.Server,
	handler http.Handler,
) {

	server := &http.Server{
		Addr:    listenAddress,
		Handler: handler,
	}

	serverConfig.HTTPServerTimeouts.ApplyToHTTPServer(server)

	if serverConfig.TLSInfo != nil {
		log.Printf("before ListenAndServeTLS serverID = %q listenAddress = %q", serverConfig.ServerID, listenAddress)

		err := server.ListenAndServeTLS(
			serverConfig.TLSInfo.CertFile,
			serverConfig.TLSInfo.KeyFile,
		)

		log.Fatalf("server.ListenAndServeTLS err = %v serverID = %q listenAddress = %q", err, serverConfig.ServerID, listenAddress)

	} else {
		log.Printf("before ListenAndServe serverID = %q listenAddress = %q", serverConfig.ServerID, listenAddress)

		err := server.ListenAndServe()

		log.Fatalf("server.ListenAndServe err = %v serverID = %q listenAddress = %q", err, serverConfig.ServerID, listenAddress)
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
