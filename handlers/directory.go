package handlers

import (
	"net/http"

	"github.com/aaronriekenberg/go-httpd/config"
)

func createDirectoryLocationHandler(
	httpPathPrefix string,
	directoryLocation config.DirectoryLocation,
) http.Handler {

	logger.Printf("createDirectoryLocationHandler httpPathPrefix = %q", httpPathPrefix)

	fileServer := http.StripPrefix(
		directoryLocation.StripPrefix,
		http.FileServer(
			http.Dir(
				directoryLocation.DirectoryPath,
			),
		),
	)

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			setCacheControlHeader(w, directoryLocation.CacheControlValue)

			fileServer.ServeHTTP(w, r)
		},
	)
}
