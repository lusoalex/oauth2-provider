package oauth2Provider

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Oauth2Flow string
type ResponseType string

type AuthorizationRequest struct {
	ClientId       ClientId
	ResponseType   ResponseType
	RedirectUri    string
	Scope          string
	State          string
	discardRequest bool
}

func (authRequest *AuthorizationRequest) invalidateRequest() {
	authRequest.discardRequest = true
}

const (
	PARAM_CLIENT_ID     = "client_id"
	PARAM_RESPONSE_TYPE = "response_type"
	PARAM_REDIRECT_URI  = "redirect_uri"
	PARAM_SCOPE         = "scope"
	PARAM_STATE         = "state"

	CODE  ResponseType = "code"
	TOKEN ResponseType = "token"

	INVALID_RESPONSE_TYPE = "Unsupported response_type"
	INVALID_REDIREC_URI   = "Unknown or invalid redirect_uri"
)

//array can't be constants...
var REQUIRED_PARAMETERS = [3]string{PARAM_CLIENT_ID, PARAM_RESPONSE_TYPE, PARAM_REDIRECT_URI}

func handleAuthorizationRequest(w http.ResponseWriter, r *http.Request) {

	var authorizationRequest AuthorizationRequest

	//Check all required parameters are well informed.
	safeInitialization(w, r, &authorizationRequest, checkRequiredParameters)

	//Initialize response_type
	safeInitialization(w, r, &authorizationRequest, initResponseType)

	//initialize client_id
	safeInitialization(w, r, &authorizationRequest, initClientId)

	//Initialize redirect_uri
	safeInitialization(w, r, &authorizationRequest, initRedirectUri)

	//If request has been discarded, then exit
	if authorizationRequest.discardRequest {
		return
	}

	//Reply with the token
	w.Header().Set(CONTENT_TYPE, CONTENT_TYPE_JSON)
	w.WriteHeader(200)
	at := "yoloooo"
	json.NewEncoder(w).Encode(Token{&at, nil})

}

func safeInitialization(w http.ResponseWriter, r *http.Request, authRequest *AuthorizationRequest, fn func(http.ResponseWriter, *http.Request, *AuthorizationRequest)) {
	//Discard request if it has been invalidated
	if authRequest.discardRequest {
		return
	} else {
		fn(w, r, authRequest)
	}
}

func checkRequiredParameters(w http.ResponseWriter, r *http.Request, authRequest *AuthorizationRequest) {
	var missingParameter []string

	for _, requiredParameter := range REQUIRED_PARAMETERS {
		if r.URL.Query().Get(requiredParameter) == "" {
			missingParameter = append(missingParameter, requiredParameter)
		}
	}

	if len(missingParameter) > 0 {
		handleError(w, errors.New("Missing parameter "+strings.Join(missingParameter, ", ")), http.StatusBadRequest)
		authRequest.invalidateRequest()
	}
}

func initResponseType(w http.ResponseWriter, r *http.Request, authRequest *AuthorizationRequest) {

	responseType := ResponseType(r.URL.Query().Get(PARAM_RESPONSE_TYPE))

	if CODE != responseType && TOKEN != responseType {
		handleError(w, errors.New(INVALID_RESPONSE_TYPE), http.StatusBadRequest)
		authRequest.invalidateRequest()
	} else {
		authRequest.ResponseType = responseType
	}
}

func initClientId(w http.ResponseWriter, r *http.Request, authRequest *AuthorizationRequest) {

	if clientId, err := findAndLoadClientSettings(r.URL.Query().Get(PARAM_CLIENT_ID)); err != nil {
		handleError(w, err, http.StatusBadRequest)
		authRequest.invalidateRequest()
	} else {
		authRequest.ClientId = *clientId
	}
}

func initRedirectUri(w http.ResponseWriter, r *http.Request, authRequest *AuthorizationRequest) {

	redirectUri := r.URL.Query().Get(PARAM_REDIRECT_URI)
	isRedirectUriValid := false

	for _, allowedRedirectUri := range authRequest.ClientId.AllowedRedirectUri {
		if redirectUri == allowedRedirectUri {
			isRedirectUriValid = true
		}
	}

	if isRedirectUriValid {
		authRequest.RedirectUri = redirectUri
	} else {
		handleError(w, errors.New(INVALID_REDIREC_URI), http.StatusBadRequest)
		authRequest.invalidateRequest()
	}
}
