package handlers

import (
	"log"
	"net/http"
	"time"

	"oauth2-provider/errors"
)

type Oauth2Handler struct{}

func (h *Oauth2Handler) ServeHTTP(_w http.ResponseWriter, req *http.Request) errors.Error {
	w := &WrappedResponseWriter{
		ResponseWriter: _w,
		StatusCode:     http.StatusOK,
	}
	defer func(path string, start time.Time) {
		log.Printf("%s %s %d %s", req.Method, path, w.StatusCode, time.Since(start))
	}(req.URL.String(), time.Now())

	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	switch head {
	case "health_check":
		return (&HealthCheckHandler{}).ServeHTTP(w, req)
	case "authorize":
		return (&AuthorizeHandler{}).ServeHTTP(w, req)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
		return nil //todo replace with an error
	}

	//router := vestigo.NewRouter()

	//Replace http.HandleFunc by vestigo compliant router.
	//router.Get("/health_check", HealthCheckHandler)
	//router.Get("/authorize", AuthorizationRequestHandler)
	//router.Post("/token", TokenRequestHandler)

	//Set custom NotFound & MethodNotAllowed handler to keep same response format...
	//vestigo.CustomNotFoundHandlerFunc(Oauth2NotFoundHandler)
	//vestigo.CustomMethodNotAllowedHandlerFunc(func(a string) func(w http.ResponseWriter, r *http.Request) {
	//	return Oauth2MethodNotAllowedHandler
	//})

	//return router
}
