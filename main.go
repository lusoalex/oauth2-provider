package main

import (
	"net/http"
	"oauth2-provider/handlers"
)

type App struct {
	handlers.Oauth2Handler
}

func main() {
	http.ListenAndServe(":8000", &App{})
}
