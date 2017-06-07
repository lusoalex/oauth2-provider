package oauth2Provider

import (
	"encoding/json"
	"net/http"
	"regexp"
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
	//Request query parameters as specified here : https://tools.ietf.org/html/rfc6749#section-4.1.1
	PARAM_CLIENT_ID     = "client_id"
	PARAM_RESPONSE_TYPE = "response_type"
	PARAM_REDIRECT_URI  = "redirect_uri"
	PARAM_SCOPE         = "scope"
	PARAM_STATE         = "state"
	//Additional query parameter as specified here : https://tools.ietf.org/html/rfc7636#section-6.1
	PARAM_CODE_VERIFIER         = "code_verifier"
	PARAM_CODE_CHALLENGE        = "code_challenge"
	PARAM_CODE_CHALLENGE_METHOD = "code_challenge_method"

	RESPONSE_TYPE_CODE  ResponseType = "code"
	RESPONSE_TYPE_TOKEN ResponseType = "token"

	CODE_CHALLENGE_METHOD_PLAIN = "plain"
	CODE_CHALLENGE_METHOD_S256  = "S256"
)

var validCodeVerifier = regexp.MustCompile("^[a-zA-Z0-9_\\-~\\.]{43,128}$")

func handleAuthorizationRequest(w http.ResponseWriter, r *http.Request) {

	var authorizationRequest AuthorizationRequest

	//Initialize response_type (see https://tools.ietf.org/html/rfc6749#section-4.1.2.1 or https://tools.ietf.org/html/rfc6749#section-4.2.2.1)
	if responseType, err := initResponseType(r); err != nil {
		handleOauth2Error(w, err)
		return
	} else {
		authorizationRequest.ResponseType = responseType
	}

	//initialize client_id
	if clientId, err := findAndLoadClientSettings(r.URL.Query().Get(PARAM_CLIENT_ID)); err != nil {
		handleOauth2Error(w, err)
		return
	} else {
		authorizationRequest.ClientId = *clientId
	}

	authorizationRequest.Scope = r.URL.Query().Get(PARAM_SCOPE)
	authorizationRequest.State = r.URL.Query().Get(PARAM_STATE)

	//Handle authorization code flow request
	switch authorizationRequest.ResponseType {
	case RESPONSE_TYPE_CODE:
		if err := handleAuthorizationCodeFlowRequest(w, r, &authorizationRequest); err != nil {
			handleOauth2Error(w, err)
			return
		}
	case RESPONSE_TYPE_TOKEN:
		if err := handleImplicitFlowRequest(w, r, &authorizationRequest); err != nil {
			handleOauth2Error(w, err)
			return
		}
	default:
		return
	}

	//Reply with the token
	w.Header().Set(CONTENT_TYPE, CONTENT_TYPE_JSON)
	w.WriteHeader(200)
	at := "yoloooo"
	json.NewEncoder(w).Encode(Token{&at, nil})

}

func initResponseType(r *http.Request) (ResponseType, *Oauth2Error) {

	responseType := ResponseType(r.URL.Query().Get(PARAM_RESPONSE_TYPE))

	if RESPONSE_TYPE_CODE != responseType && RESPONSE_TYPE_TOKEN != responseType {
		return "", NewResponseTypeError()
	}
	return responseType, nil
}

/*
 *  On the authorization code flow, the redirect_uri is required : https://tools.ietf.org/html/rfc6749#section-4.1.1
 *  But on implicit flow, it is not mandatory as specified here : https://tools.ietf.org/html/rfc6749#section-4.2.1
 *  In such case we must ensure that the request come's from an allowed client uri https://tools.ietf.org/html/rfc6749#section-3.1.2
 */
func initRedirectUri(r *http.Request, allowedRedirectUris []string, isImplicit bool) (string, *Oauth2Error) {

	//If redirect_uri is not informed and current request is oauth2 implicit flow, then we get it from the settings.
	if redirectUri := r.URL.Query().Get(PARAM_REDIRECT_URI); redirectUri == "" && isImplicit && len(allowedRedirectUris) == 1 {
		return allowedRedirectUris[0], nil
	} else {
		//check that the provided redirect_uri is well informed into the client settings.
		for _, allowedRedirectUri := range allowedRedirectUris {
			if redirectUri == allowedRedirectUri {
				return redirectUri, nil
			}
		}
	}
	//No matching redirect_uri found, return an error.
	return "", NewRedirectUriError(isImplicit)
}

/**
 * Even if PKCE (https://tools.ietf.org/html/rfc7636) is not forced, if code_challenge is informed, we will apply it.
 */
func handleAuthorizationCodeFlowRequest(w http.ResponseWriter, r *http.Request, authRequest *AuthorizationRequest) *Oauth2Error {

	//Initialize redirect_uri (required query parameter)
	if redirectUri, err := initRedirectUri(r, authRequest.ClientId.AllowedRedirectUri, false); err != nil {
		return err
	} else {
		authRequest.RedirectUri = redirectUri
	}

	//Get code_challenge, and if client_id settings require use of PKCE, return an error if not respected.
	codeChallenge := r.URL.Query().Get(PARAM_CODE_CHALLENGE)
	if codeChallenge == "" && authRequest.ClientId.ForceUseOfPKCE {
		return NewCodeChallengeError()
	}

	codeChallengeMethod := r.URL.Query().Get(PARAM_CODE_CHALLENGE_METHOD)

	//If code_challenge_method is specified, then the value must be plain or S256
	if codeChallengeMethod != "" && codeChallengeMethod != CODE_CHALLENGE_METHOD_PLAIN && codeChallengeMethod != CODE_CHALLENGE_METHOD_S256 {
		return NewCodeChallengeMethodError()
	}

	//If the code_challenge_method is not specified, but there's a code_challenge informed, so we use plain as default
	//For more details, see : https://tools.ietf.org/html/rfc7636#section-4.3
	if codeChallenge != "" && codeChallengeMethod == "" {
		codeChallengeMethod = CODE_CHALLENGE_METHOD_PLAIN
	}

	return nil
}

func handleImplicitFlowRequest(w http.ResponseWriter, r *http.Request, authRequest *AuthorizationRequest) *Oauth2Error {

	//Initialize redirect_uri (optional query parameter)
	if redirectUri, err := initRedirectUri(r, authRequest.ClientId.AllowedRedirectUri, true); err != nil {
		return err
	} else {
		authRequest.RedirectUri = redirectUri
	}

	return nil
}

/*
 * As specified in specs https://tools.ietf.org/html/rfc7636#section-4.1
 * code_verifier = high-entropy cryptographic random STRING using the
 * unreserved characters [A-Z] / [a-z] / [0-9] / "-" / "." / "_" / "~"
 * length must be between 43 to 128 characters
 */
func valideCodeVerifier(codeVerifier string) *Oauth2Error {
	m := validCodeVerifier.FindStringSubmatch(codeVerifier)
	if m == nil {
		return newCodeVerifierFormatError()
	}
	return nil
}
