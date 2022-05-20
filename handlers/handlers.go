package handlers

import (
	"net/http"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/aaronriekenberg/go-httpd/logging"
)

var logger = logging.GetLogger()

func CreateServerLocationsHandler(
	locations []config.Location,
	customResponseHeaders *config.CustomResponseHeaders,
) http.Handler {

	var handler http.Handler = newLocationListHandler(
		locations,
	)

	handler = createCustomResponseHeadersHandler(
		customResponseHeaders,
		handler,
	)

	return handler
}
