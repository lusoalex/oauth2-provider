package handlers

import (
	"net/http"

	"encoding/json"
)

type HealthCheckHandler struct{}

// TODO In the future we could report back on the status of our DB, or our cache
// TODO (e.g. Redis) by performing a simple PING, and include them in the response.
func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	switch head {
	case "":
		switch req.Method {
		case "GET":
			if bytes, err := json.Marshal(&struct {
				Alive bool `json:"alive"`
			}{Alive: true}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(bytes)
			}
		}
	}
}
