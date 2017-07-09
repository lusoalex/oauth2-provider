package handlers

import (
	"net/http"

	oauth2_errors "oauth2-provider/errors"
)

type Oauth2ErrorHandler struct {
	oauth2_errors.Error
}

func (h *Oauth2ErrorHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) error {

}
