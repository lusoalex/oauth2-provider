package main

import (
	"net/http"
	"./handlers"
)

type App struct {
	handlers.Oauth2Handler
}



func main() {
	http.ListenAndServe(":8000", &App{})
}
