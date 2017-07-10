package handlers

import (
	"net/http"
	"oauth2-provider/response"
)

type MainHandler struct{}

func (h MainHandler) Handle(w http.ResponseWriter, req *http.Request) (*response.HTTPResponse, error) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	switch head {
	case "health_check":
		return (&HealthCheckHandler{}).Handle(w, req)
	case "authorize":
		return (&AuthorizeHandler{}).Handle(w, req)
	}

	return nil, NotFound
}
