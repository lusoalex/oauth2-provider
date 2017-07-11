package handlers

import (
	"bytes"
	"net/http"
	"net/url"
	"oauth2-provider/client"
	"oauth2-provider/constants"
	oauth2_errors "oauth2-provider/errors"
	"oauth2-provider/models"
	"oauth2-provider/response"
	"oauth2-provider/user"
	"strings"
)

type AuthorizeHandler struct {
	MainHandler
}

func (*AuthorizeHandler) Handle(w http.ResponseWriter, r *http.Request) {

	var authorizationRequest models.AuthorizationRequest

	//initialize client_id
	if clientId, clientIdErr := client.GetClientInformations(r.URL.Query().Get(constants.PARAM_CLIENT_ID)); clientIdErr != nil {
		//TODO return response.BadRequest(response.NewJsonResponse(client.INVALID_CLIENT_ID))
	} else {
		authorizationRequest.ClientId = *clientId
	}

	authorizationRequest.ResponseType = models.ResponseType(r.URL.Query().Get(constants.PARAM_RESPONSE_TYPE))
	authorizationRequest.Scope = r.URL.Query().Get(constants.PARAM_SCOPE)
	authorizationRequest.State = r.URL.Query().Get(constants.PARAM_STATE)

	//Initialize redirect_uri
	if redirectUri := initRedirectUri(r, authorizationRequest.ClientId.AllowedRedirectUri); redirectUri == "" {
		//TODO return response.BadRequest(response.NewJsonResponse(oauth2_errors.InvalidRedirectUri))
	} else {
		authorizationRequest.RedirectUri = redirectUri
	}

	var err error
	//Handle authorization code flow request
	switch authorizationRequest.ResponseType {
	case models.RESPONSE_TYPE_CODE:
		err = handleAuthorizationCodeFlowRequest(w, r, &authorizationRequest)
	case models.RESPONSE_TYPE_TOKEN:
		err = handleImplicitFlowRequest(w, r, &authorizationRequest)
	default:
		//err = 500
	}

	if err != nil {
		switch err.(type) {
		case models.BadRequest:
			w.WriteHeader(http.StatusBadRequest)
			//w.Write(json(err))
		case models.ForbiddenRequest:
			w.WriteHeader(http.StatusForbidden)
			//w.Write(json(err))
		default:
			w.WriteHeader(500)
			//w.Write(???)
		}
	}

	return
}

/**
 * Even if PKCE (https://tools.ietf.org/html/rfc7636) is not forced, if code_challenge is informed, we will apply it.
 */
func handleAuthorizationCodeFlowRequest(w http.ResponseWriter, r *http.Request, authRequest *models.AuthorizationRequest) error {

	if !isGrantTypeAllowed(models.GRANT_TYPE_AUTHORIZATION_CODE, authRequest.ClientId.AllowedGrantType) {
		return models.BadRequest(oauth2_errors.UnauthorizedClient(models.GRANT_TYPE_AUTHORIZATION_CODE))
	}

	//Get code_challenge, and if client_id settings require use of PKCE, return an error if not respected.
	codeChallenge := r.URL.Query().Get(constants.PARAM_CODE_CHALLENGE)
	if codeChallenge == "" && authRequest.ClientId.ForceUseOfPKCE {
		return oauth2_errors.MissingCodeChallenge
	}

	codeChallengeMethod := models.CodeChallengeMethod(r.URL.Query().Get(constants.PARAM_CODE_CHALLENGE_METHOD))
	//If code_challenge_method is specified, then the value must be plain or S256
	if codeChallengeMethod != "" && codeChallengeMethod != models.CODE_CHALLENGE_METHOD_PLAIN && codeChallengeMethod != models.CODE_CHALLENGE_METHOD_S256 {
		return oauth2_errors.InvalidCodeChallenge
	}

	//If the code_challenge_method is not specified, but there's a code_challenge informed, so we use plain as default
	//For more details, see : https://tools.ietf.org/html/rfc7636#section-4.3
	if codeChallenge != "" && codeChallengeMethod == "" {
		codeChallengeMethod = models.CODE_CHALLENGE_METHOD_PLAIN
	}

	authRequest.CodeChallenge = codeChallenge
	authRequest.CodeChallengeMethod = codeChallengeMethod

	//TODO display login form

	/*
		//generate code
		authRequest.Code = uuid.New().String()
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go storeCode(&wg, authRequest)

		//build redirect uri
		uri, _ := url.Parse(authRequest.RedirectUri)
		query := uri.Query()
		query.Add(constants.PARAM_CODE, authRequest.Code)

		uri.RawQuery = query.Encode()
		http.Redirect(w, r, uri.String(), http.StatusFound)

		wg.Wait() //wait code have been stored to return the response

		return nil, errors.New("Not implemented yet")
	*/
	return nil
}

func handleImplicitFlowRequest(w http.ResponseWriter, r *http.Request, authRequest *models.AuthorizationRequest) error {

	if !isGrantTypeAllowed(models.GRANT_TYPE_IMPLICIT, authRequest.ClientId.AllowedGrantType) {
		return models.BadRequest(oauth2_errors.UnauthorizedClient(models.GRANT_TYPE_IMPLICIT))
	}

	if form, user := HandleLoginPage(r); user == nil {
		//w.XXX
		response.NewHTTPResponse(form, "text/html")
	} else {
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
func initRedirectUri(r *http.Request, allowedRedirectUris []string) string {

	//TODO : check wildcard uri (be as restrictive as possible)
	//TODO : implement optional redirect_uri

	// If redirect_uri is not informed and current request is oauth2 implicit flow, then we get it from the settings.
	if redirectUri := r.URL.Query().Get(constants.PARAM_REDIRECT_URI); redirectUri != "" {
		//check that the provided redirect_uri is well informed into the client settings.
		for _, allowedRedirectUri := range allowedRedirectUris {
			if redirectUri == allowedRedirectUri {
				return redirectUri
			}
		}
	}

	//No matching redirect_uri found, return an error.
	return ""
}

func isGrantTypeAllowed(grantType models.GrantType, allowedGrantType []models.GrantType) bool {
	for _, allowedGrant := range allowedGrantType {
		if grantType == allowedGrant {
			return true
		}
	}
	return false
}

/*
func storeCode(wg *sync.WaitGroup, authRequest *AuthorizationRequest) error {

	defer wg.Done()

	var value bytes.Buffer
	enc := gob.NewEncoder(&value)
	if err := enc.Encode(authRequest); err != nil {
		//TODO how to manage erors in goroutine??
		return err
	}

	getKeyValueStore().Set(authRequest.Code, value, 20*time.Second)

	return nil
}
*/

func HandleLoginPage(r *http.Request) ([]byte, *user.User) {

	r.ParseForm()
	if login, password := r.Form.Get("login"), r.Form.Get("password"); r.Method == "POST" && login != "" && password != "" {
		if user, ok := user.MatchingCredentials(login, password); ok == true {
			return nil, user
		}
	}

	var w bytes.Buffer
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

	return w.Bytes(), nil
}
