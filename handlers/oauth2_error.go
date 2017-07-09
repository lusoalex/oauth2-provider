package handlers

import (
	"oauth2-provider/errors"
	"net/http"
)

type Oauth2ErrorHandler struct {
	errors.Oauth2Error
}

func (h *Oauth2ErrorHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}
