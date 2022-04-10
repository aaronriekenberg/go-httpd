package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/kr/pretty"
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

			redirectURL := strings.ReplaceAll(redirectLocation.RedirectURL, "$REQUEST_PATH", r.URL.Path)

			http.Redirect(w, r, redirectURL, redirectLocation.ResponseStatus)

		},
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
	log.Printf("requestURLPath = %q", requestURLPath)

	for _, locationHandler := range locationListHandler.locationHandlers {

		match := strings.HasPrefix(requestURLPath, locationHandler.httpPathPrefix)

		log.Printf("match = %v", match)

		if match {
			locationHandler.httpHandler.ServeHTTP(w, r)
			return
		}
	}

	log.Printf("got to end of locationHandlers with no match, requestURLPath = %q", requestURLPath)
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func CreateLocationsHandler(
	locations []config.Location,
) http.Handler {

	handler := &locationListHandler{}

	for _, location := range locations {
		log.Printf("location:\n%# v", pretty.Formatter(location))

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

		default:
			log.Fatalf("invalid location config: \n%# v", pretty.Formatter(location))
		}
	}

	return handler
}
