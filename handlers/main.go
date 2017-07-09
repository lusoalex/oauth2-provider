package handlers

import (
	"net/http"
	"log"
)

type MainHandler struct{
	Oauth2ProviderHandler
}

func (h *MainHandler) Handle(w http.ResponseWriter, req *http.Request) error {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	log.Printf("Main : current head : %s", head)

	switch head {
	case "health_check":
		return (&HealthCheckHandler{}).Handle(w, req)
	case "authorize":
		return (&AuthorizeHandler{}).Handle(w, req)
	}

	return NotFound

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
