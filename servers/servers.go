package servers

import (
	"log"
	"net/http"
	"os"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/kr/pretty"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/aaronriekenberg/go-httpd/handlers"
)

func runServer(
	serverConfig config.Server,
) {
	handler := handlers.CreateLocationsHandler(serverConfig.Locations)

	if serverConfig.LogRequests {
		handler = gorillaHandlers.CombinedLoggingHandler(os.Stdout, handler)
	}

	server := &http.Server{
		Addr:    serverConfig.ListenAddress,
		Handler: handler,
	}

	if serverConfig.TLSInfo != nil {
		log.Printf("before ListenAndServeTLS listenAddress = %q", serverConfig.ListenAddress)

		err := server.ListenAndServeTLS(
			serverConfig.TLSInfo.CertFile,
			serverConfig.TLSInfo.KeyFile,
		)

		log.Fatalf("server.ListenAndServeTLS err = %v", err)

	} else {
		log.Printf("before ListenAndServe listenAddress = %q", serverConfig.ListenAddress)

		err := server.ListenAndServe()

		log.Fatalf("server.ListenAndServe err = %v", err)
	}
}

func StartServers(
	servers []config.Server,
) {
	log.Printf("begin StartServers")

	for _, serverConfig := range servers {
		log.Printf("serverConfig:\n%# v", pretty.Formatter(serverConfig))

		go runServer(serverConfig)
	}

	log.Printf("end StartServers")
}
