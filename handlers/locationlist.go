package handlers

import (
	"net/http"

	"github.com/aaronriekenberg/go-httpd/config"
)

type locationListHandler struct {
	locationHandlers      []*locationHandler
	customResponseHeaders *config.CustomResponseHeaders
}

func (locationListHandler *locationListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var matchingLocationHandler *locationHandler

	for _, locationHandler := range locationListHandler.locationHandlers {

		if locationHandler.matches(r) {
			matchingLocationHandler = locationHandler
			break
		}
	}

	locationListHandler.customResponseHeaders.ApplyToResponse(w)

	if matchingLocationHandler == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	matchingLocationHandler.ServeHTTP(w, r)
}
