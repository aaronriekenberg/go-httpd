package handlers

import (
	"net/http"

	"github.com/aaronriekenberg/go-httpd/config"
)

func NewServerLocationsHandler(
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
