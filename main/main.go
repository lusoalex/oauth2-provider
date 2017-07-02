package main

import (
	"github.com/lusoalex/oauth2-provider"
)

func main() {
	oauth2Provider.Serve(&oauth2Provider.Oauth2ServerOptions{
		Port: "8000",
	})
}
