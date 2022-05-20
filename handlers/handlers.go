package handlers

import (
	"net/http"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/aaronriekenberg/go-httpd/logging"
)

var logger = logging.GetLogger()

func CreateLocationsHandler(
	locations []config.Location,
	customResponseHeaders *config.CustomResponseHeaders,
) http.Handler {

	var handler locationListHandler

	for _, locationConfig := range locations {

		handler.locationHandlers = append(
			handler.locationHandlers,
			newLocationHandler(locationConfig),
		)

	}

	handler.customResponseHeaders = customResponseHeaders

	return &handler
}
