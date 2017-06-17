package oauth2Provider

import (
	"net/http"

	"github.com/husobee/vestigo"
)

func Oauth2Handler() http.Handler {

	router := vestigo.NewRouter()

	//Replace http.HandleFunc by vestigo compliant router.
	router.Get("/health_check", HealthCheckHandler)
	router.Get("/authorize", AuthorizationRequestHandler)
	router.Post("/token", TokenRequestHandler)

	//Set custom NotFound & MethodNotAllowed handler to keep same response format...
	vestigo.CustomNotFoundHandlerFunc(Oauth2NotFoundHandler)
	vestigo.CustomMethodNotAllowedHandlerFunc(func(a string) func(w http.ResponseWriter, r *http.Request) {
		return Oauth2MethodNotAllowedHandler
	})

	return router
}
