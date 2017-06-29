package oauth2Provider

import (
	"net/http"

	"github.com/husobee/vestigo"
)

var kvs KeyValueStore

type Oauth2ServerOptions struct {
	Port string
	Kvs  KeyValueStore
}

func setKeyValueStore(new KeyValueStore) {
	kvs = new
}

func getKeyValueStore() KeyValueStore {
	return kvs
}

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

func Serve(options *Oauth2ServerOptions) {
	setKeyValueStore(options.Kvs)
	http.ListenAndServe(":"+options.Port, Oauth2Handler())
}
