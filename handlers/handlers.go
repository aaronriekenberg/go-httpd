package handlers

import (
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/kr/pretty"
	"github.com/lpar/gzipped/v2"
	"github.com/yookoala/gofast"

	"github.com/aaronriekenberg/go-httpd/config"
)

type locationHandler struct {
	httpPathPrefix string
	httpHandler    http.Handler
}

func addCacheControlHeader(
	w http.ResponseWriter,
	cacheControlValue string,
) {
	if len(cacheControlValue) > 0 {
		w.Header().Add("cache-control", cacheControlValue)
	}
}

func createBlockedLocationHandler(
	httpPathPrefix string,
	blockedLocation config.BlockedLocation,
) http.Handler {

	log.Printf("createBlockedLocationHandler httpPathPrefix = %q", httpPathPrefix)

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(blockedLocation.ResponseStatus)
		},
	)
}

func createDirectoryLocationHandler(
	httpPathPrefix string,
	directoryLocation config.DirectoryLocation,
) http.Handler {

	log.Printf("createDirectoryLocationHandler httpPathPrefix = %q", httpPathPrefix)

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
			addCacheControlHeader(w, directoryLocation.CacheControlValue)

			fileServer.ServeHTTP(w, r)
		},
	)
}

func createCompressedDirectoryLocationHandler(
	httpPathPrefix string,
	compressedDirectoryLocation config.CompressedDirectoryLocation,
) http.Handler {

	log.Printf("createCompressedDirectoryLocationHandler httpPathPrefix = %q", httpPathPrefix)

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
			addCacheControlHeader(w, compressedDirectoryLocation.CacheControlValue)

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

func createRedirectLocationHandler(
	httpPathPrefix string,
	redirectLocation config.RedirectLocation,
) http.Handler {

	log.Printf("createRedirectLocationHandler httpPathPrefix = %q", httpPathPrefix)

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			if r.URL == nil {
				http.Error(w, "request url is nil in redirect handler", http.StatusInternalServerError)
				return
			}

			redirectURL := redirectLocation.RedirectURL
			redirectURL = strings.ReplaceAll(redirectURL, "$HTTP_HOST", r.Host)
			redirectURL = strings.ReplaceAll(redirectURL, "$REQUEST_PATH", r.URL.Path)

			http.Redirect(w, r, redirectURL, redirectLocation.ResponseStatus)

		},
	)
}

func createFastCGILocationHandler(
	httpPathPrefix string,
	fastCGILocation config.FastCGILocation,
) http.Handler {

	log.Printf("createFastCGILocationHandler httpPathPrefix = %q", httpPathPrefix)

	sessionHandler := gofast.Chain(
		gofast.BasicParamsMap, // maps common CGI parameters
		gofast.MapHeader,      // maps header fields into HTTP_* parameters
	)(gofast.BasicSession)

	connectionFactory := gofast.SimpleConnFactory("unix", fastCGILocation.UnixSocketPath)

	// XXX make parameters configurable?
	connectionPool := gofast.NewClientPool(
		gofast.SimpleClientFactory(connectionFactory),
		10,             // buffer size for pre-created client-connection
		10*time.Second, // life span of a client before expire
	)

	return gofast.NewHandler(
		sessionHandler,
		connectionPool.CreateClient,
	)
}

type locationListHandler struct {
	locationHandlers []*locationHandler
}

func (locationListHandler *locationListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestURLPath := r.URL.Path

	var handler http.Handler

	for _, locationHandler := range locationListHandler.locationHandlers {

		match := strings.HasPrefix(requestURLPath, locationHandler.httpPathPrefix)

		if match {
			handler = locationHandler.httpHandler
			break
		}
	}

	if handler == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	handler.ServeHTTP(w, r)
}

func createHandlerForLocation(
	locationConfig config.Location,
) http.Handler {

	var locationHandler http.Handler

	switch {
	case locationConfig.BlockedLocation != nil:
		locationHandler = createBlockedLocationHandler(
			locationConfig.HttpPathPrefix,
			*locationConfig.BlockedLocation,
		)

	case locationConfig.DirectoryLocation != nil:
		locationHandler = createDirectoryLocationHandler(
			locationConfig.HttpPathPrefix,
			*locationConfig.DirectoryLocation,
		)

	case locationConfig.CompressedDirectoryLocation != nil:
		locationHandler = createCompressedDirectoryLocationHandler(
			locationConfig.HttpPathPrefix,
			*locationConfig.CompressedDirectoryLocation,
		)

	case locationConfig.RedirectLocation != nil:
		locationHandler = createRedirectLocationHandler(
			locationConfig.HttpPathPrefix,
			*locationConfig.RedirectLocation,
		)

	case locationConfig.FastCGILocation != nil:
		locationHandler = createFastCGILocationHandler(
			locationConfig.HttpPathPrefix,
			*locationConfig.FastCGILocation,
		)

	}

	if locationHandler == nil {
		log.Fatalf("invalid location config: \n%# v", pretty.Formatter(locationConfig))
	}

	return locationHandler
}

func CreateLocationsHandler(
	locations []config.Location,
) http.Handler {

	handler := &locationListHandler{}

	for _, locationConfig := range locations {

		handler.locationHandlers = append(
			handler.locationHandlers,
			&locationHandler{
				httpPathPrefix: locationConfig.HttpPathPrefix,
				httpHandler:    createHandlerForLocation(locationConfig),
			},
		)

	}

	return handler
}
