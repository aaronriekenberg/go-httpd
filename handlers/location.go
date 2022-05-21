package handlers

import (
	"net/http"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/kr/pretty"
)

type locationHandler struct {
	requestMatcher
	http.Handler
}

func newLocationHandler(locationConfig config.Location) *locationHandler {
	return &locationHandler{
		requestMatcher: newRequestMatcher(locationConfig.HttpPathPrefix),
		Handler:        createHTTPHandlerForLocation(locationConfig),
	}
}

func createHTTPHandlerForLocation(
	locationConfig config.Location,
) http.Handler {

	var locationHandler http.Handler

	switch {
	case locationConfig.BlockedLocation != nil:
		locationHandler = newBlockedLocationHandler(
			locationConfig.HttpPathPrefix,
			*locationConfig.BlockedLocation,
		)

	case locationConfig.DirectoryLocation != nil:
		locationHandler = newDirectoryLocationHandler(
			locationConfig.HttpPathPrefix,
			*locationConfig.DirectoryLocation,
		)

	case locationConfig.CompressedDirectoryLocation != nil:
		locationHandler = newCompressedDirectoryLocationHandler(
			locationConfig.HttpPathPrefix,
			*locationConfig.CompressedDirectoryLocation,
		)

	case locationConfig.RedirectLocation != nil:
		locationHandler = newRedirectLocationHandler(
			locationConfig.HttpPathPrefix,
			*locationConfig.RedirectLocation,
		)

	case locationConfig.FastCGILocation != nil:
		locationHandler = newFastCGILocationHandler(
			locationConfig.HttpPathPrefix,
			*locationConfig.FastCGILocation,
		)

	}

	if locationHandler == nil {
		logger.Fatalf("invalid location config: \n%# v", pretty.Formatter(locationConfig))
	}

	locationHandler = newResponseHeadersHandler(
		locationConfig.ResponseHeaders,
		locationHandler,
	)

	return locationHandler
}
