package handlers

import (
	"net/http"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/aaronriekenberg/go-httpd/logging"
)

var logger = logging.GetLogger()

func setCacheControlHeader(
	w http.ResponseWriter,
	cacheControlValue string,
) {
	if len(cacheControlValue) > 0 {
		w.Header().Set("cache-control", cacheControlValue)
	}
}

func CreateLocationsHandler(
	locations []config.Location,
) http.Handler {

	handler := &locationListHandler{}

	for _, locationConfig := range locations {

		handler.locationHandlers = append(
			handler.locationHandlers,
			newLocationHandler(locationConfig),
		)

	}

	return handler
}
