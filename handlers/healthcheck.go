package handlers

import (
	"io"
	"net/http"
)

type HealthCheckHandler struct {

}

func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	switch head {
	case "":
		switch req.Method {
		case "GET":
			// A very simple health check.
			w.Header().Set(CONTENT_TYPE, CONTENT_TYPE_JSON)
			w.WriteHeader(http.StatusOK)
			// TODO In the future we could report back on the status of our DB, or our cache
			// TODO (e.g. Redis) by performing a simple PING, and include them in the response.
			io.WriteString(w, `{"alive": true}`)
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}