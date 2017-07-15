package handlers

import (
	"net/http"
	"net/url"
	"oauth2-provider/client"
	"oauth2-provider/constants"
	oauth2_errors "oauth2-provider/errors"
	"oauth2-provider/models"
	"oauth2-provider/user"
	"strings"
)

type AuthorizeHandler struct {
	Oauth2Handler
}

func (a *AuthorizeHandler) ServeHttp(w http.ResponseWriter, r *http.Request) {
	a.handleOauth2Request("authorize", []string{"GET", "POST"}, w, r, a.handleAuthorizationRequest)
}

func (a *AuthorizeHandler) handleAuthorizationRequest(w http.ResponseWriter, r *http.Request) error {
	var authorizationRequest models.AuthorizationRequest

	//Initialize client_id
	client, err := client.GetClientInformations(r.URL.Query().Get(constants.PARAM_CLIENT_ID))
	if err != nil {
		return err
	}
	authorizationRequest.ClientId = client.ClientId

	//Initialize redirect_uri
	if authorizationRequest.RedirectUri, err = initRedirectUri(r, client.AllowedRedirectUri); err != nil {
		return err
	}

	//Initialize Score, State and ResponseType
	authorizationRequest.Scope = r.URL.Query().Get(constants.PARAM_SCOPE)
	authorizationRequest.State = r.URL.Query().Get(constants.PARAM_STATE)
	authorizationRequest.ResponseType = models.ResponseType(r.URL.Query().Get(constants.PARAM_RESPONSE_TYPE))

	//Handle request according to the ResponseType
	switch authorizationRequest.ResponseType {
	case models.RESPONSE_TYPE_CODE:
		return a.handleAuthorizationCodeFlowRequest(w, r, &authorizationRequest, client)
	case models.RESPONSE_TYPE_TOKEN:
		return a.handleImplicitFlowRequest(w, r, &authorizationRequest, client)
	default:
		return &models.BadRequest{Oauth2Error: &models.Oauth2Error{
			Reason:           "unsupported_response_type",
			ErrorDescription: "Missing, unsupported or malformed required response_type parameter.",
			ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-3.1.1",
		}}
	}

	return nil
}

/**
 * Even if PKCE (https://tools.ietf.org/html/rfc7636) is not forced, if code_challenge is informed, we will apply it.
 */
func (a *AuthorizeHandler) handleAuthorizationCodeFlowRequest(w http.ResponseWriter, r *http.Request, authRequest *models.AuthorizationRequest, client *models.ClientId) error {

	if !isGrantTypeAllowed(models.GRANT_TYPE_AUTHORIZATION_CODE, client.AllowedGrantType) {
		return oauth2_errors.UnauthorizedClient(models.GRANT_TYPE_AUTHORIZATION_CODE)
	}

	//Get code_challenge, and if client_id settings require use of PKCE, return an error if not respected.
	codeChallenge := r.URL.Query().Get(constants.PARAM_CODE_CHALLENGE)
	if codeChallenge == "" && client.ForceUseOfPKCE {
		return oauth2_errors.InvalidRequest("Missing required code_challenger parameter.", "https://tools.ietf.org/html/rfc7636#section-4.4.1")
	}

	codeChallengeMethod := models.CodeChallengeMethod(r.URL.Query().Get(constants.PARAM_CODE_CHALLENGE_METHOD))
	//If code_challenge_method is specified, then the value must be plain or S256
	if codeChallengeMethod != "" && codeChallengeMethod != models.CODE_CHALLENGE_METHOD_PLAIN && codeChallengeMethod != models.CODE_CHALLENGE_METHOD_S256 {
		return oauth2_errors.InvalidRequest("Invalid code_challange_method parameter", "https://tools.ietf.org/html/rfc7636#section-4.3")
	}

	//If the code_challenge_method is not specified, but there's a code_challenge informed, so we use plain as default
	//For more details, see : https://tools.ietf.org/html/rfc7636#section-4.3
	if codeChallenge != "" && codeChallengeMethod == "" {
		codeChallengeMethod = models.CODE_CHALLENGE_METHOD_PLAIN
	}

	authRequest.CodeChallenge = codeChallenge
	authRequest.CodeChallengeMethod = codeChallengeMethod

	if user := HandleLoginPage(w, r); user != nil {

		//generate code
		code := a.Code(authRequest)

		//build redirect uri
		uri, _ := url.Parse(authRequest.RedirectUri)
		query := uri.Query()
		query.Add(constants.PARAM_CODE, code)

		uri.RawQuery = query.Encode()
		http.Redirect(w, r, uri.String(), http.StatusFound)
	}

	return nil
}

func (a *AuthorizeHandler) handleImplicitFlowRequest(w http.ResponseWriter, r *http.Request, authRequest *models.AuthorizationRequest, client *models.ClientId) error {

	if !isGrantTypeAllowed(models.GRANT_TYPE_IMPLICIT, client.AllowedGrantType) {
		return oauth2_errors.UnauthorizedClient(models.GRANT_TYPE_IMPLICIT)
	}

	if user := HandleLoginPage(w, r); user != nil {
		//TODO generate true jwt
		accessToken := "token." + user.Firstname + "." + user.Name
		//build redirect uri
		uri, _ := url.Parse(authRequest.RedirectUri)
		query := uri.Query()
		query.Add(constants.PARAM_ACCESS_TOKEN, accessToken)
		query.Add(constants.PARAM_TOKEN_TYPE, string(constants.TOKEN_TYPE_BEARER))
		//TODO add expires_in
		//TODO add scope if different
		if authRequest.State != "" {
			query.Add(constants.PARAM_STATE, authRequest.State)
		}
		uri.RawQuery = query.Encode()
		location := strings.Replace(uri.String(), "?", "#", 1)

		http.Redirect(w, r, location, http.StatusFound)
	}

	return nil
}

/*
 * On the authorization code flow, the redirect_uri is optional : https://tools.ietf.org/html/rfc6749#section-4.1.1
 * And on implicit flow too as specified here : https://tools.ietf.org/html/rfc6749#section-4.2.1
 * In such case we must ensure that the request come's from an allowed client uri https://tools.ietf.org/html/rfc6749#section-3.1.2
 * The endpoint URI MUST NOT include a fragment component. TODO: add additional check in/out values
 */
func initRedirectUri(r *http.Request, allowedRedirectUris []string) (string, error) {

	//TODO : check wildcard uri (be as restrictive as possible)

	// If redirect_uri is not informed and current request is oauth2 implicit flow, then we get it from the settings.
	if redirectUri := r.URL.Query().Get(constants.PARAM_REDIRECT_URI); redirectUri == "" && len(allowedRedirectUris) == 1 {
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
	return "", oauth2_errors.InvalidRequest(
		"Missing, invalid, or mismatching redirect_uri parameter.",
		"https://tools.ietf.org/html/rfc6749#section-3.1.2",
	)
}

func isGrantTypeAllowed(grantType models.GrantType, allowedGrantType []models.GrantType) bool {
	for _, allowedGrant := range allowedGrantType {
		if grantType == allowedGrant {
			return true
		}
	}
	return false
}

func HandleLoginPage(w http.ResponseWriter, r *http.Request) *user.User {

	r.ParseForm()
	if login, password := r.Form.Get("login"), r.Form.Get("password"); r.Method == "POST" && login != "" && password != "" {
		if user, ok := user.MatchingCredentials(login, password); ok == true {
			return user
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "text/html")

	w.Write([]byte("<!DOCTYPE html>"))
	w.Write([]byte("<html><head>"))
	w.Write([]byte("<meta charset=\"UTF-8\">"))
	w.Write([]byte("<title>Login Page</title>"))
	w.Write([]byte("</head><body>"))
	w.Write([]byte("<div><form method=\"post\" >"))
	w.Write([]byte("Login<br><input type=\"text\" name=\"login\" placeholder=\"Login\" required=\"required\"/><br><br>"))
	w.Write([]byte("Password<br><input type=\"password\" name=\"password\" placeholder=\"Password\" required=\"required\"/><br><br>"))
	w.Write([]byte("<button type=\"submit\">LOGIN</button>"))
	w.Write([]byte("</form></div></body></html>"))

	return nil
}
