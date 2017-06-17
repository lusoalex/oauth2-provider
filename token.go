package oauth2Provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type GrantType string
type TokenType string
type Token struct {
	AccessToken  *string    `json:"access_token"`
	TokenType    *TokenType `json:"token_type,omitempty"`
	ExpireIn     *int       `json:"expires_in,omitempty"`
	RefreshToken *string    `json:"refresh_token,omitempty"`
}

const (
	GRANT_TYPE_REFRESH_TOKEN      GrantType = "refresh_token"
	GRANT_TYPE_AUTHORIZATION_CODE GrantType = "authorization_code"
	GRANT_TYPE_PASSWORD           GrantType = "password"
	GRANT_TYPE_CLIENT_CREDENTIALS GrantType = "client_credentials"

	TOKEN_TYPE_BEARER TokenType = "Bearer"
)

var validCodeVerifier = regexp.MustCompile("^[a-zA-Z0-9_-~.]{43,128}$")

func TokenRequestHandler(w http.ResponseWriter, r *http.Request) {

	grant_type := GrantType(r.URL.Query().Get(PARAM_GRANT_TYPE))

	var err *Oauth2Error
	//Handle request following the grant_type
	switch grant_type {
	case GRANT_TYPE_AUTHORIZATION_CODE:
		fmt.Println("TODO: token authorization code")
		err = handleAuthorizationCodeTokenRequest(w, r)
	case GRANT_TYPE_CLIENT_CREDENTIALS:
		fmt.Println("TODO: token client crendentials")
	case GRANT_TYPE_PASSWORD:
		fmt.Println("TODO: token password")
	case GRANT_TYPE_REFRESH_TOKEN:
		fmt.Println("TODO: token refresh token")
	default:
		err = NewGrantTypeError()
	}

	if err != nil {
		handleOauth2Error(w, err)
		return
	}

	//Reply with the token
	w.Header().Set(CONTENT_TYPE, CONTENT_TYPE_JSON)
	w.WriteHeader(200)
	at := "yoloooo"
	rt := "god bless you"
	json.NewEncoder(w).Encode(Token{AccessToken: &at, RefreshToken: &rt})

}

func handleAuthorizationCodeTokenRequest(w http.ResponseWriter, r *http.Request) *Oauth2Error {

	//initialize client_id
	clientId, err := findAndLoadClientSettings(r.URL.Query().Get(PARAM_CLIENT_ID))
	if err != nil {
		return err
	}

	//Manage code_verifier if need
	if codeVerifier := r.URL.Query().Get(PARAM_CODE_VERIFIER); codeVerifier != "" || clientId.ForceUseOfPKCE {
		if err := validateCodeVerifier(codeVerifier); err != nil {
			return err
		}
	}

	return nil
}

/*
 * As specified in specs https://tools.ietf.org/html/rfc7636#section-4.1
 * code_verifier = high-entropy cryptographic random STRING using the
 * unreserved characters [A-Z] / [a-z] / [0-9] / "-" / "." / "_" / "~"
 * length must be between 43 to 128 characters
 */
func validateCodeVerifier(codeVerifier string) *Oauth2Error {

	if m := validCodeVerifier.FindStringSubmatch(codeVerifier); m == nil {
		return NewCodeVerifierFormatError()
	}

	//TODO validate code verifier corresponds to code_challenge given on the authorize request
	if err := NewCodeVerifierError(); err != nil {
		return err
	}

	return nil
}
