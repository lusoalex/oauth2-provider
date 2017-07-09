package handlers

import (
	oauth2_errors "oauth2-provider/errors"
	"net/http"
)

type Oauth2ErrorHandler struct {
	oauth2_errors.Error
}

func (h *Oauth2ErrorHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) error {

}
