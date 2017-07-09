package handlers

import (
	"net/http"

	"log"
	"oauth2-provider/response"
)

type HealthCheckHandler struct {
	MainHandler
}

// TODO In the future we could report back on the status of our DB, or our cache
// TODO (e.g. Redis) by performing a simple PING, and include them in the response.
func (h *HealthCheckHandler) Handle(w http.ResponseWriter, req *http.Request) (response.Response, error) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	log.Printf("HealthCheck : Current head : %s\n", head)
	switch head {
	case "":
		switch req.Method {
		case "GET":
			return &response.JsonResponse{
				Content: &struct {
					Alive bool
				}{
					Alive: true,
				},
			}, nil
		}
	}
	return nil, NotFound
}
