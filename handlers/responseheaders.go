package handlers

import (
	"net/http"

	"github.com/aaronriekenberg/go-httpd/config"
)

func newResponseHeadersHandler(
	responseHeaders *config.ResponseHeaders,
	nextHandler http.Handler,
) http.Handler {

	if (responseHeaders == nil) || (len(*responseHeaders) == 0) {
		return nextHandler
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			for key, value := range *responseHeaders {
				w.Header().Set(key, value)
			}

			nextHandler.ServeHTTP(w, r)
		},
	)
}
