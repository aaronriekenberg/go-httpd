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
	serverID string,
	listenAddress string,
	tlsInfo *config.TLSInfo,
	handler http.Handler,
) {

	server := &http.Server{
		Addr:    listenAddress,
		Handler: handler,
	}

	if tlsInfo != nil {
		log.Printf("before ListenAndServeTLS serverID = %q listenAddress = %q", serverID, listenAddress)

		err := server.ListenAndServeTLS(
			tlsInfo.CertFile,
			tlsInfo.KeyFile,
		)

		log.Fatalf("server.ListenAndServeTLS err = %v serverID = %q listenAddress = %q", err, serverID, listenAddress)

	} else {
		log.Printf("before ListenAndServe serverID = %q listenAddress = %q", serverID, listenAddress)

		err := server.ListenAndServe()

		log.Fatalf("server.ListenAndServe err = %v serverID = %q listenAddress = %q", err, serverID, listenAddress)
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
				serverConfig.ServerID,
				listenAddress,
				serverConfig.TLSInfo,
				handler,
			)
		}
	}

	log.Printf("end StartServers")
}
