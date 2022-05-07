package handlers

import "net/http"

type locationListHandler struct {
	locationHandlers []*locationHandler
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
