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
	ClientId     ClientId
	ResponseType ResponseType
	RedirectUri  string
	Scope        string
	State        string
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
	INVALID_REDIRECT_URI  = "Unknown or invalid redirect_uri"
)

func handleAuthorizationRequest(w http.ResponseWriter, r *http.Request) {

	var authorizationRequest AuthorizationRequest

	//Check all required parameters are well informed.
	if err := checkRequiredParameters(r); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	//Initialize response_type
	if responseType, err := initResponseType(r); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	} else {
		authorizationRequest.ResponseType = responseType
	}

	//initialize client_id
	if clientId, err := findAndLoadClientSettings(r.URL.Query().Get(PARAM_CLIENT_ID)); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	} else {
		authorizationRequest.ClientId = *clientId
	}

	//Initialize redirect_uri
	if redirectUri, err := initRedirectUri(r, authorizationRequest.ClientId.AllowedRedirectUri); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	} else {
		authorizationRequest.RedirectUri = redirectUri
	}

	authorizationRequest.Scope = r.URL.Query().Get(PARAM_SCOPE)
	authorizationRequest.State = r.URL.Query().Get(PARAM_STATE)

	//Reply with the token
	w.Header().Set(CONTENT_TYPE, CONTENT_TYPE_JSON)
	w.WriteHeader(200)
	at := "yoloooo"
	json.NewEncoder(w).Encode(Token{&at, nil})

}

func checkRequiredParameters(r *http.Request) error {
	var REQUIRED_PARAMETERS = [3]string{PARAM_CLIENT_ID, PARAM_RESPONSE_TYPE, PARAM_REDIRECT_URI}
	var missingParameter []string

	for _, requiredParameter := range REQUIRED_PARAMETERS {
		if r.URL.Query().Get(requiredParameter) == "" {
			missingParameter = append(missingParameter, requiredParameter)
		}
	}

	if len(missingParameter) > 0 {
		return errors.New("Missing parameter " + strings.Join(missingParameter, ", "))
	}
	return nil
}

func initResponseType(r *http.Request) (ResponseType, error) {

	responseType := ResponseType(r.URL.Query().Get(PARAM_RESPONSE_TYPE))

	if CODE != responseType && TOKEN != responseType {
		return "", errors.New(INVALID_RESPONSE_TYPE)
	}
	return responseType, nil
}

func initRedirectUri(r *http.Request, allowedRedirectUris []string) (string, error) {

	redirectUri := r.URL.Query().Get(PARAM_REDIRECT_URI)
	isRedirectUriValid := false

	for _, allowedRedirectUri := range allowedRedirectUris {
		if redirectUri == allowedRedirectUri {
			isRedirectUriValid = true
		}
	}

	if isRedirectUriValid {
		return redirectUri, nil
	} else {
		return "", errors.New(INVALID_REDIRECT_URI)
	}
}
