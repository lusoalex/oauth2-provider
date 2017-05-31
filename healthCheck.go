package oauth2Provider

import (
	"io"
	"net/http"
)

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set(CONTENT_TYPE, CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusOK)
	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}
