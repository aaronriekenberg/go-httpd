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

	return http.StripPrefix(
		directoryLocation.StripPrefix,
		http.FileServer(
			http.Dir(
				directoryLocation.DirectoryPath,
			),
		),
	)
}
