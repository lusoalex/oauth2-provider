package handlers

import (
	"net/http"

	"oauth2-provider/response"
)

type HealthCheckHandler struct {
	MainHandler
}

// TODO In the future we could report back on the status of our DB, or our cache
// TODO (e.g. Redis) by performing a simple PING, and include them in the response.
func (h *HealthCheckHandler) Handle(w http.ResponseWriter, req *http.Request) (*response.HTTPResponse, error) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	switch head {
	case "":
		switch req.Method {
		case "GET":
			return response.OK(response.NewJsonResponse(
				&struct {
					Alive bool `json:"alive"`
				}{
					Alive: true,
				},
			))
		}
	}
	return nil, NotFound
}
