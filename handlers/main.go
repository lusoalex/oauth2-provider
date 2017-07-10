package handlers

import (
	"log"
	"net/http"
	"oauth2-provider/response"
)

type MainHandler struct{}

func (h MainHandler) Handle(w http.ResponseWriter, req *http.Request) (response.Response, error) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	log.Printf("Main : current head : %s", head)

	switch head {
	case "health_check":
		return (&HealthCheckHandler{}).Handle(w, req)
	case "authorize":
		return (&AuthorizeHandler{}).Handle(w, req)
	}

	return nil, NotFound
}
