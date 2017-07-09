package main

import (
	"net/http"
	"oauth2-provider/handlers"
	"log"
)

type App struct {
	handlers.MainHandler
}

func main() {
	if serveErr := http.ListenAndServe(":8000", &App{}) ; serveErr != nil {
		log.Fatalln(serveErr)
	}
}
