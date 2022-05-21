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

	sessionHandler := gofast.Chain(
		gofast.BasicParamsMap, // maps common CGI parameters
		gofast.MapHeader,      // maps header fields into HTTP_* parameters
	)(gofast.BasicSession)

	connectionFactory := gofast.SimpleConnFactory("unix", fastCGILocation.UnixSocketPath)

	connectionPoolSize := uint(10)
	if fastCGILocation.ConnectionPoolSize != nil {
		connectionPoolSize = uint(*fastCGILocation.ConnectionPoolSize)
	}

	connectionPoolLifetimeDuration := 10 * time.Second
	if fastCGILocation.ConnectionPoolLifetimeMilliseconds != nil {
		connectionPoolLifetimeDuration = time.Duration(*fastCGILocation.ConnectionPoolLifetimeMilliseconds) * time.Millisecond
	}

	logger.Printf(
		"newFastCGILocationHandler httpPathPrefix = %q connectionPoolSize = %v connectionPoolLifetimeDuration = %v",
		httpPathPrefix,
		connectionPoolSize,
		connectionPoolLifetimeDuration,
	)

	connectionPool := gofast.NewClientPool(
		gofast.SimpleClientFactory(connectionFactory),
		connectionPoolSize,             // buffer size for pre-created client-connection
		connectionPoolLifetimeDuration, // life span of a client before expire
	)

	return gofast.NewHandler(
		sessionHandler,
		connectionPool.CreateClient,
	)
}
