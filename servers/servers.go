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
	listenAddress string,
	tlsInfo *config.TLSInfo,
	handler http.Handler,
) {

	server := &http.Server{
		Addr:    listenAddress,
		Handler: handler,
	}

	if tlsInfo != nil {
		log.Printf("before ListenAndServeTLS listenAddress = %q", listenAddress)

		err := server.ListenAndServeTLS(
			tlsInfo.CertFile,
			tlsInfo.KeyFile,
		)

		log.Fatalf("server.ListenAndServeTLS err = %v", err)

	} else {
		log.Printf("before ListenAndServe listenAddress = %q", listenAddress)

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

		handler := handlers.CreateLocationsHandler(serverConfig.Locations)

		if serverConfig.LogRequests {
			handler = gorillaHandlers.CombinedLoggingHandler(os.Stdout, handler)
		}

		for _, listenAddress := range serverConfig.ListenAddressList {
			go runServer(
				listenAddress,
				serverConfig.TLSInfo,
				handler,
			)
		}
	}

	log.Printf("end StartServers")
}
