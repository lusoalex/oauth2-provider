package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	oauth2_errors "oauth2-provider/errors"
	"oauth2-provider/constants"
	"oauth2-provider/client"
	"oauth2-provider/utils"
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

	grant_type := GrantType(r.URL.Query().Get(constants.PARAM_GRANT_TYPE))

	var err *oauth2_errors.Oauth2Error
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
		err = oauth2_errors.GrantTypeError
	}

	if err != nil {
		err.Handle(w)
		return
	}

	//Reply with the token
	w.Header().Set(constants.CONTENT_TYPE, constants.CONTENT_TYPE_JSON)
	w.WriteHeader(200)
	at := "yoloooo"
	rt := "god bless you"
	json.NewEncoder(w).Encode(Token{AccessToken: &at, RefreshToken: &rt})

}

func handleAuthorizationCodeTokenRequest(w http.ResponseWriter, r *http.Request) *oauth2_errors.Oauth2Error {
	var authRequest AuthorizationRequest

	//initialize client_id
	clientId, err := client.FindAndLoadClientSettings(r.URL.Query().Get(constants.PARAM_CLIENT_ID))
	if err != nil {
		return err
	}

	//Manage code_verifier if need
	if codeVerifier := r.URL.Query().Get(constants.PARAM_CODE_VERIFIER); codeVerifier != "" || clientId.ForceUseOfPKCE {
		if err := validateCodeVerifier(codeVerifier); err != nil {
			return err
		}
	}

	//TODO retrieve request from key value store and check all parameters are matching...
	//code := r.URL.Query().Get(constants.PARAM_CODE)
	//byteCode, _ := base64.RawURLEncoding.DecodeString(code)
	//byteRequest, _ := utils.KVS.Get(authRequest)

	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)
	//buf.Write(byteRequest)
	dec.Decode(&authRequest)
	code := utils.KVS.Get(authRequest)

	return nil
}

/*
 * As specified in specs https://tools.ietf.org/html/rfc7636#section-4.1
 * code_verifier = high-entropy cryptographic random STRING using the
 * unreserved characters [A-Z] / [a-z] / [0-9] / "-" / "." / "_" / "~"
 * length must be between 43 to 128 characters
 */
func validateCodeVerifier(codeVerifier string) *oauth2_errors.Oauth2Error {

	if m := validCodeVerifier.FindStringSubmatch(codeVerifier); m == nil {
		return oauth2_errors.CodeVerifierFormatError
	}

	//TODO validate code verifier corresponds to code_challenge given on the authorize request
	if err := oauth2_errors.CodeVerifierError; err != nil {
		return err
	}

	return nil
}