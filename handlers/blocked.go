package handlers

import (
	"net/http"

	"github.com/aaronriekenberg/go-httpd/config"
)

func createBlockedLocationHandler(
	httpPathPrefix string,
	blockedLocation config.BlockedLocation,
) http.Handler {

	logger.Printf("createBlockedLocationHandler httpPathPrefix = %q", httpPathPrefix)

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(blockedLocation.ResponseStatus)
		},
	)
}
