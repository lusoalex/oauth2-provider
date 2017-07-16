package main

import (
	"log"
	"net/http"
	"oauth2-provider/handlers"
	"oauth2-provider/settings"
)

func main() {
	handler := &handlers.MainHandler{settings.DefaultOauth2ProviderSettings()}
	//handler := &handlers.HealthCheckHandler{app.DefaultOauth2ProviderSettings()} //--> No sense as "health_check" path is not expected anymore
	//handler := &handlers.AuthorizeHandler{app.DefaultOauth2ProviderSettings()}   //--> No sense as "authorize" path is not expected anymore
	//handler := &handlers.TokenHandler{app.DefaultOauth2ProviderSettings()}       //--> No sense as "token" path is not expected anymore

	if err := http.ListenAndServe(":8000", handler); err != nil {
		log.Fatalln(err)
	}
}
