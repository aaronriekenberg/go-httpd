package handlers

import (
	"net/http"
	"strings"
)

type requestMatcher struct {
	httpPathPrefix string
}

func newRequestMatcher(
	httpPathPrefix string,
) requestMatcher {
	return requestMatcher{
		httpPathPrefix: httpPathPrefix,
	}
}

func (requestMatcher *requestMatcher) matches(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, requestMatcher.httpPathPrefix)
}
