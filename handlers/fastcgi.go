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

	var clientFactory gofast.ClientFactory

	logger.Printf(
		"newFastCGILocationHandler httpPathPrefix = %q connectionPool = %+v",
		httpPathPrefix,
		fastCGILocation.ConnectionPool,
	)

	if !fastCGILocation.ConnectionPool.Enabled {
		clientFactory = gofast.SimpleClientFactory(connectionFactory)
	} else {
		connectionPoolSize := uint(10)
		if fastCGILocation.ConnectionPool.Size != nil {
			connectionPoolSize = uint(*fastCGILocation.ConnectionPool.Size)
		}

		connectionPoolLifetimeDuration := 10 * time.Second
		if fastCGILocation.ConnectionPool.LifetimeMilliseconds != nil {
			connectionPoolLifetimeDuration = time.Duration(*fastCGILocation.ConnectionPool.LifetimeMilliseconds) * time.Millisecond
		}

		connectionPool := gofast.NewClientPool(
			gofast.SimpleClientFactory(connectionFactory),
			connectionPoolSize,             // buffer size for pre-created client-connection
			connectionPoolLifetimeDuration, // life span of a client before expire
		)

		clientFactory = connectionPool.CreateClient
	}

	return gofast.NewHandler(
		sessionHandler,
		clientFactory,
	)
}
