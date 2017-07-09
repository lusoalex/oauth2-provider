package main

import (
	"net/http"
	"oauth2-provider/handlers"
	"log"
)

type App struct {
	handlers.CommonHandler
}

func main() {
	if serveErr := http.ListenAndServe(":8000", &App{
		handlers.CommonHandler{
			Handler: handlers.MainHandler{},
		},
	}) ; serveErr != nil {
		log.Fatalln(serveErr)
	}
}
