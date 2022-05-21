package handlers

import (
	"net/http"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/aaronriekenberg/go-httpd/logging"
)

var logger = logging.GetLogger()

func CreateServerLocationsHandler(
	locations []config.Location,
	serverResponseHeaders *config.ResponseHeaders,
) http.Handler {

	var handler http.Handler = newLocationListHandler(
		locations,
	)

	handler = newResponseHeadersHandler(
		serverResponseHeaders,
		handler,
	)

	return handler
}
