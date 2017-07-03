package main

import (
	"net/http"
	"./handlers"
)

func main() {
	http.ListenAndServe(":8000", &handlers.Oauth2Handler{})
}
