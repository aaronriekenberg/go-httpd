package handlers

import (
	"net/http"
	"path"
	"strings"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/lpar/gzipped/v2"
)

func createCompressedDirectoryLocationHandler(
	httpPathPrefix string,
	compressedDirectoryLocation config.CompressedDirectoryLocation,
) http.Handler {

	logger.Printf("createCompressedDirectoryLocationHandler httpPathPrefix = %q", httpPathPrefix)

	fileServer := http.StripPrefix(
		compressedDirectoryLocation.StripPrefix,
		gzipped.FileServer(
			gzipped.Dir(
				compressedDirectoryLocation.DirectoryPath,
			),
		),
	)

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			setCacheControlHeader(w, compressedDirectoryLocation.CacheControlValue)

			// Unlike http.FileServer, gzipped.FileServer does not serve
			// index.html for directory requests by default.
			// See withIndexHTML example:
			// https://github.com/lpar/gzipped/blob/trunk/README.md
			if strings.HasSuffix(r.URL.Path, "/") || len(r.URL.Path) == 0 {
				newpath := path.Join(r.URL.Path, "index.html")
				r.URL.Path = newpath
			}

			fileServer.ServeHTTP(w, r)
		},
	)

}
