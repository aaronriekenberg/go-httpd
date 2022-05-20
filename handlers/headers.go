package handlers

import (
	"net/http"

	"github.com/aaronriekenberg/go-httpd/config"
)

func setCacheControlHeader(
	w http.ResponseWriter,
	cacheControlValue string,
) {
	const cacheControlKey = "cache-control"

	if len(cacheControlValue) > 0 {
		w.Header().Set(cacheControlKey, cacheControlValue)
	}
}

func createCustomResponseHeadersHandler(
	customResponseHeaders *config.CustomResponseHeaders,
	nextHandler http.Handler,
) http.Handler {

	if customResponseHeaders == nil {
		return nextHandler
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			for key, value := range *customResponseHeaders {
				w.Header().Set(key, value)
			}

			nextHandler.ServeHTTP(w, r)
		},
	)
}
