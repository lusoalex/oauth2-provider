package handlers

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"regexp"

	"fmt"
	"oauth2-provider/client"
	"oauth2-provider/constants"
	"oauth2-provider/models"
)

type TokenHandler struct{}
type TokenType string
type Token struct {
	AccessToken  *string    `json:"access_token"`
	TokenType    *TokenType `json:"token_type,omitempty"`
	ExpireIn     *int       `json:"expires_in,omitempty"`
	RefreshToken *string    `json:"refresh_token,omitempty"`
}

const (
	TOKEN_TYPE_BEARER TokenType = "Bearer"
)

var validCodeVerifier = regexp.MustCompile("^[a-zA-Z0-9_-~.]{43,128}$")

func (*TokenHandler) ServeHttp(w http.ResponseWriter, r *http.Request) {
	HandleOauth2Request(w, r, handleTokenRequest)
}

func handleTokenRequest(w http.ResponseWriter, r *http.Request) error {
	grant_type := models.GrantType(r.URL.Query().Get(constants.PARAM_GRANT_TYPE))

	//Handle request following the grant_type
	switch grant_type {
	case models.GRANT_TYPE_AUTHORIZATION_CODE:
		fmt.Println("TODO: token authorization code")
		return handleAuthorizationCodeTokenRequest(w, r)
	case models.GRANT_TYPE_CLIENT_CREDENTIALS:
		fmt.Println("TODO: token client crendentials")
	case models.GRANT_TYPE_RESSOURCE_OWNER_PASSWORD_CREDENTIALS:
		fmt.Println("TODO: token password")
	case models.GRANT_TYPE_REFRESH_TOKEN:
		fmt.Println("TODO: token refresh token")
	default:
		return models.BadRequest{Oauth2Error: &models.Oauth2Error{
			Reason:           "unsupported_grant_type",
			ErrorDescription: "Missing, unsupported or malformed required grant_type parameter.",
			ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-5.2",
		}}
	}

	/* TODO move this...
	//Reply with the token
	w.Header().Set(constants.CONTENT_TYPE, constants.CONTENT_TYPE_JSON)
	w.WriteHeader(200)
	at := "yoloooo"
	rt := "god bless you"
	json.NewEncoder(w).Encode(Token{AccessToken: &at, RefreshToken: &rt})
	*/
}

func handleAuthorizationCodeTokenRequest(w http.ResponseWriter, r *http.Request) error {
	var authRequest models.AuthorizationRequest

	//initialize client_id
	clientId, err := client.GetClientInformations(r.URL.Query().Get(constants.PARAM_CLIENT_ID))
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
	//code := utils.KVS.Get(authRequest)

	return nil
}

/*
 * As specified in specs https://tools.ietf.org/html/rfc7636#section-4.1
 * code_verifier = high-entropy cryptographic random STRING using the
 * unreserved characters [A-Z] / [a-z] / [0-9] / "-" / "." / "_" / "~"
 * length must be between 43 to 128 characters
 */
func validateCodeVerifier(codeVerifier string) error {

	//if m := validCodeVerifier.FindStringSubmatch(codeVerifier); m == nil {
	//	return oauth2_errors.InvalidCodeVerifier
	//}

	//TODO validate code verifier corresponds to code_challenge given on the authorize request
	/*if err := oauth2_errors.CodeVerifierError; err != nil {
		return err
	}*/

	return nil
}
