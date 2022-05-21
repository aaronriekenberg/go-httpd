package handlers

import (
	"net/http"
	"time"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/yookoala/gofast"
)

func newFastCGILocationHandler(
	httpPathPrefix string,
	fastCGILocation config.FastCGILocation,
) http.Handler {

	logger.Printf("newFastCGILocationHandler httpPathPrefix = %q", httpPathPrefix)

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
