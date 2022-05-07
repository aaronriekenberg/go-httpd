package handlers

import "net/http"

func setCacheControlHeader(
	w http.ResponseWriter,
	cacheControlValue string,
) {
	const cacheControlKey = "cache-control"

	if len(cacheControlValue) > 0 {
		w.Header().Set(cacheControlKey, cacheControlValue)
	}
}
