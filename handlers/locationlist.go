package handlers

import (
	"net/http"

	"github.com/aaronriekenberg/go-httpd/config"
)

type locationListHandler struct {
	locationHandlers []*locationHandler
}

func newLocationListHandler(
	locations []config.Location,
) http.Handler {

	handler := &locationListHandler{
		locationHandlers: make([]*locationHandler, 0, len(locations)),
	}

	for _, locationConfig := range locations {

		handler.locationHandlers = append(
			handler.locationHandlers,
			newLocationHandler(locationConfig),
		)

	}

	return handler
}

func (locationListHandler *locationListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var matchingLocationHandler *locationHandler

	for _, locationHandler := range locationListHandler.locationHandlers {

		if locationHandler.matches(r) {
			matchingLocationHandler = locationHandler
			break
		}
	}

	if matchingLocationHandler == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	matchingLocationHandler.ServeHTTP(w, r)
}
