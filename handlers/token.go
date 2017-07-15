package handlers

import (
	"net/http"
	"regexp"

	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"oauth2-provider/client"
	"oauth2-provider/constants"
	"oauth2-provider/errors"
	"oauth2-provider/models"
)

type TokenHandler struct {
	Oauth2Handler
}

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

func (t *TokenHandler) ServeHttp(w http.ResponseWriter, r *http.Request) {
	t.handleOauth2Request("token", []string("POST"), w, r, t.handleTokenRequest)
}

func (t *TokenHandler) handleTokenRequest(w http.ResponseWriter, r *http.Request) error {

	var token *Token
	var err error

	//Handle request following the grant_type
	switch models.GrantType(r.URL.Query().Get(constants.PARAM_GRANT_TYPE)) {
	case models.GRANT_TYPE_AUTHORIZATION_CODE:
		token, err = t.handleAuthorizationCodeTokenRequest(w, r)
	case models.GRANT_TYPE_CLIENT_CREDENTIALS:
		//TODO: token client crendentials
		token, err = nil, nil
	case models.GRANT_TYPE_RESSOURCE_OWNER_PASSWORD_CREDENTIALS:
		//TODO: token password
		token, err = nil, nil
	case models.GRANT_TYPE_REFRESH_TOKEN:
		//"TODO: token refresh token
		token, err = nil, nil
	default:
		token, err = nil, models.BadRequest{Oauth2Error: &models.Oauth2Error{
			Reason:           "unsupported_grant_type",
			ErrorDescription: "Missing, unsupported or malformed required grant_type parameter.",
			ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-5.2",
		}}
	}

	if err != nil {
		return err
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(token)
	return nil
}

func (t *TokenHandler) handleAuthorizationCodeTokenRequest(w http.ResponseWriter, r *http.Request) (*Token, error) {

	//initialize client_id
	clientId, err := client.GetClientInformations(r.URL.Query().Get(constants.PARAM_CLIENT_ID))
	if err != nil {
		return nil, err
	}
	//TODO check client authentication and it corresponds to current clientId

	//Get the original Authorize request from the code parameter
	authRequest, err := t.extractAuthorizationRequestFromCode(r)
	if err != nil {
		return nil, err
	}

	//Check client_id is the same between authorize and token request
	if authRequest.ClientId != clientId.ClientId {
		return nil, errors.InvalidGrant("Client id does not match the id of the client to whom the authorization code was issued.")
	}

	//Manage code_verifier if need
	if codeVerifier := r.URL.Query().Get(constants.PARAM_CODE_VERIFIER); codeVerifier != "" || clientId.ForceUseOfPKCE {
		if err := validateCodeVerifier(codeVerifier, authRequest); err != nil {
			return nil, err
		}
	}

	at := "yoloooo"
	rt := "god bless you"

	return &Token{AccessToken: &at, RefreshToken: &rt}, nil
}

func (t *TokenHandler) extractAuthorizationRequestFromCode(r *http.Request) (*models.AuthorizationRequest, error) {
	//Get code from the token request
	if code := r.URL.Query().Get(constants.PARAM_CODE); code == "" {
		return nil, errors.InvalidRequest(
			"Missing, invalid, or mismatching code parameter.",
			"https://tools.ietf.org/html/rfc6749#section-5.2",
		)
	} else {
		//Get the associated authorize request, as code must be used only once, revoke it.
		if ar, ok := t.Revoke(code); !ok {
			//Unable to retrieve /authorize Request from the code
			return nil, errors.InvalidGrant("Authorization code is invalid or expired.")
		} else {
			return ar, nil
		}
	}
}

/*
 * As specified in specs https://tools.ietf.org/html/rfc7636#section-4.1
 * code_verifier = high-entropy cryptographic random STRING using the
 * unreserved characters [A-Z] / [a-z] / [0-9] / "-" / "." / "_" / "~"
 * length must be between 43 to 128 characters
 */
func validateCodeVerifier(codeVerifier string, authRequest *models.AuthorizationRequest) error {

	if m := validCodeVerifier.FindStringSubmatch(codeVerifier); m == nil {
		return errors.InvalidRequest(
			"Missing or malformed code_verifier parameter",
			"https://tools.ietf.org/html/rfc7636#section-4.1",
		)
	}

	//Validate the code_verifier according to code challenge method provided in the authorize request.
	var matchingCodeVerifier bool
	switch authRequest.CodeChallengeMethod {
	case models.CODE_CHALLENGE_METHOD_PLAIN:
		matchingCodeVerifier = authRequest.CodeChallenge == codeVerifier
	case models.CODE_CHALLENGE_METHOD_S256:
		//BASE64URL-ENCODE(SHA256(ASCII(code_verifier))) == code_challenge
		matchingCodeVerifier = authRequest.CodeChallenge == base64.RawURLEncoding.EncodeToString(sha256.Sum256([]byte(codeVerifier))[:])
	}

	//Return error if code verifier is not valid
	if !matchingCodeVerifier {
		err := errors.InvalidGrant("Invalid code_verifier parameter")
		err.ErrorUri = "https://tools.ietf.org/html/rfc7636#section-4.6" //override ErrorUri as PKCE come's from another RFC.
		return err
	}

	return nil
}
