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

	logger.Printf(
		"newFastCGILocationHandler httpPathPrefix = %q fastCGILocation = %+v",
		httpPathPrefix,
		fastCGILocation,
	)

	sessionHandler := gofast.Chain(
		gofast.BasicParamsMap, // maps common CGI parameters
		gofast.MapHeader,      // maps header fields into HTTP_* parameters
	)(gofast.BasicSession)

	connectionFactory := gofast.SimpleConnFactory(fastCGILocation.Network, fastCGILocation.Address)

	var clientFactory gofast.ClientFactory

	if fastCGILocation.ConnectionPool == nil {
		clientFactory = gofast.SimpleClientFactory(connectionFactory)
	} else {
		connectionPoolSize := uint(fastCGILocation.ConnectionPool.Size)

		connectionPoolLifetimeDuration := time.Duration(fastCGILocation.ConnectionPool.LifetimeMilliseconds) * time.Millisecond

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
