package handlers

import (
	"net/http"
	"strings"

	"github.com/aaronriekenberg/go-httpd/config"
)

func newRedirectLocationHandler(
	httpPathPrefix string,
	redirectLocation config.RedirectLocation,
) http.Handler {

	logger.Printf("newRedirectLocationHandler httpPathPrefix = %q", httpPathPrefix)

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			redirectURL := redirectLocation.RedirectURL
			redirectURL = strings.ReplaceAll(redirectURL, "$HTTP_HOST", r.Host)
			redirectURL = strings.ReplaceAll(redirectURL, "$REQUEST_PATH", r.URL.Path)

			http.Redirect(w, r, redirectURL, redirectLocation.ResponseStatus)

		},
	)
}
