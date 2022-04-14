package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/kr/pretty"
	"github.com/yookoala/gofast"

	"github.com/aaronriekenberg/go-httpd/config"
)

type locationHandler struct {
	httpPathPrefix string
	httpHandler    http.Handler
}

func createBlockedLocationHandler(
	httpPathPrefix string,
	blockedLocation config.BlockedLocation,
) locationHandler {

	log.Printf("createBlockedLocationHandler httpPathPrefix = %q", httpPathPrefix)

	return locationHandler{
		httpPathPrefix: httpPathPrefix,
		httpHandler: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(blockedLocation.ResponseStatus)
			},
		),
	}
}

func createDirectoryLocationHandler(
	httpPathPrefix string,
	directoryLocation config.DirectoryLocation,
) locationHandler {

	log.Printf("createDirectoryLocationHandler httpPathPrefix = %q", httpPathPrefix)

	return locationHandler{
		httpPathPrefix: httpPathPrefix,
		httpHandler: http.StripPrefix(
			directoryLocation.StripPrefix,
			http.FileServer(
				http.Dir(
					directoryLocation.DirectoryPath,
				),
			),
		),
	}
}

func createRedirectLocationHandler(
	httpPathPrefix string,
	redirectLocation config.RedirectLocation,
) locationHandler {

	log.Printf("createRedirectLocationHandler httpPathPrefix = %q", httpPathPrefix)

	handler := http.HandlerFunc(
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

	return locationHandler{
		httpPathPrefix: httpPathPrefix,
		httpHandler:    handler,
	}
}

func createFastCGILocationHandler(
	httpPathPrefix string,
	fastCGILocation config.FastCGILocation,
) locationHandler {

	log.Printf("createFastCGILocationHandler httpPathPrefix = %q", httpPathPrefix)

	sessionHandler := gofast.Chain(
		gofast.BasicParamsMap, // maps common CGI parameters
		gofast.MapHeader,      // maps header fields into HTTP_* parameters
	)(gofast.BasicSession)

	connectionFactory := gofast.SimpleConnFactory("unix", fastCGILocation.UnixSocketPath)

	handler := gofast.NewHandler(
		sessionHandler,
		gofast.SimpleClientFactory(connectionFactory),
	)

	return locationHandler{
		httpPathPrefix: httpPathPrefix,
		httpHandler:    handler,
	}
}

type locationListHandler struct {
	locationHandlers []locationHandler
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

func CreateLocationsHandler(
	locations []config.Location,
) http.Handler {

	handler := &locationListHandler{}

	for _, location := range locations {

		switch {
		case location.BlockedLocation != nil:
			handler.locationHandlers = append(
				handler.locationHandlers,
				createBlockedLocationHandler(
					location.HttpPathPrefix,
					*location.BlockedLocation,
				),
			)

		case location.DirectoryLocation != nil:
			handler.locationHandlers = append(
				handler.locationHandlers,
				createDirectoryLocationHandler(
					location.HttpPathPrefix,
					*location.DirectoryLocation,
				),
			)

		case location.RedirectLocation != nil:
			handler.locationHandlers = append(
				handler.locationHandlers,
				createRedirectLocationHandler(
					location.HttpPathPrefix,
					*location.RedirectLocation,
				),
			)

		case location.FastCGILocation != nil:
			handler.locationHandlers = append(
				handler.locationHandlers,
				createFastCGILocationHandler(
					location.HttpPathPrefix,
					*location.FastCGILocation,
				),
			)

		default:
			log.Fatalf("invalid location config: \n%# v", pretty.Formatter(location))
		}
	}

	return handler
}
