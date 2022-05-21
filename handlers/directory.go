package handlers

import (
	"net/http"

	"github.com/aaronriekenberg/go-httpd/config"
)

func newDirectoryLocationHandler(
	httpPathPrefix string,
	directoryLocation config.DirectoryLocation,
) http.Handler {

	logger.Printf("newDirectoryLocationHandler httpPathPrefix = %q", httpPathPrefix)

	return http.StripPrefix(
		directoryLocation.StripPrefix,
		http.FileServer(
			http.Dir(
				directoryLocation.DirectoryPath,
			),
		),
	)
}
