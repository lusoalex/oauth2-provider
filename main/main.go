package main

import (
	"net/http"

	"github.com/lusoalex/oauth2-provider"
)

func main() {
	http.ListenAndServe(":8000", oauth2Provider.Oauth2Handler())
}
