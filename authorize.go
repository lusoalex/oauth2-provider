package oauth2Provider

import (
	"errors"
	"net/http"
	"strings"
	"fmt"
	"encoding/json"
)

type Oauth2Flow string
type ResponseType string

type ClientId struct {
	ClientId           string
	AllowedRedirectUri []string
	AllowedOauth2Flow  []Oauth2Flow
}

type AuthorizationRequest struct {
	ClientId     ClientId
	ResponseType ResponseType
	RedirectUri  string
	Scope        string
	state        string
}

const (
	PARAM_CLIENT_ID     = "client_id"
	PARAM_RESPONSE_TYPE = "response_type"
	PARAM_REDIRECT_URI  = "redirect_uri"
	PARAM_SCOPE         = "scope"
	PARAM_STATE         = "state"

	CODE ResponseType   = "code"
	TOKEN ResponseType  = "token"
)

//array can't be constants...
var REQUIRED_PARAMETERS = [3]string {PARAM_CLIENT_ID, PARAM_RESPONSE_TYPE, PARAM_REDIRECT_URI}



func handleAuthorizationRequest(w http.ResponseWriter, r *http.Request) {

	fmt.Println("handling authorization request")

	for k,v := range r.URL.Query() {
		fmt.Print(k,v)
	}

	if err := checkRequiredParameters(r) ; err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}
	at := "yoloooo"
	json.NewEncoder(w).Encode(Token{&at,nil})
}

func checkRequiredParameters(r *http.Request) error {
	var missingParameter []string

	for _,requiredParameter := range REQUIRED_PARAMETERS {
		if r.URL.Query().Get(requiredParameter) == "" {
			missingParameter = append(missingParameter,requiredParameter)
		}
	}

	if len(missingParameter) > 0 {
		return errors.New("Missing parameter "+strings.Join(missingParameter, ", "))
	}
	return nil
}

