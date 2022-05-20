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

	handler := &locationListHandler{
		locationHandlers:      make([]*locationHandler, 0, len(locations)),
		customResponseHeaders: customResponseHeaders,
	}

	for _, locationConfig := range locations {

		handler.locationHandlers = append(
			handler.locationHandlers,
			newLocationHandler(locationConfig),
		)

	}

	return handler
}
