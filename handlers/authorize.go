package handlers

import (
	"net/http"
	"oauth2-provider/client"
	"oauth2-provider/constants"
	"oauth2-provider/models"
	"oauth2-provider/response"
)

type Oauth2Flow string
type AuthorizeHandler struct {
}

func (*AuthorizeHandler) Handle(w http.ResponseWriter, r *http.Request) (response.Response, error) {

	var authorizationRequest models.AuthorizationRequest

	//initialize client_id
	if clientId, clientIdErr := client.FindAndLoadClientSettings(r.URL.Query().Get(constants.PARAM_CLIENT_ID)); clientIdErr != nil {
		return nil, clientIdErr
	} else {
		authorizationRequest.ClientId = *clientId
	}

	authorizationRequest.ResponseType = models.ResponseType(r.URL.Query().Get(constants.PARAM_RESPONSE_TYPE))
	authorizationRequest.Scope = r.URL.Query().Get(constants.PARAM_SCOPE)
	authorizationRequest.State = r.URL.Query().Get(constants.PARAM_STATE)

	//Handle authorization code flow request
	/*switch authorizationRequest.ResponseType {
	case models.RESPONSE_TYPE_CODE:
		return handleAuthorizationCodeFlowRequest(w, r, &authorizationRequest)
	case models.RESPONSE_TYPE_TOKEN:
		return handleImplicitFlowRequest(w, r, &authorizationRequest)
	default:
		return nil, errors.New("Invalid response type")
	}*/

	return nil, NotFound
}

/**
 * Even if PKCE (https://tools.ietf.org/html/rfc7636) is not forced, if code_challenge is informed, we will apply it.
 */
func handleAuthorizationCodeFlowRequest(w http.ResponseWriter, r *http.Request, authRequest *models.AuthorizationRequest) error {
/*
	//Initialize redirect_uri (required query parameter)
	if redirectUri, err := initRedirectUri(r, authRequest.ClientId.AllowedRedirectUri, false); err != nil {
		return err
	} else {
		authRequest.RedirectUri = redirectUri
	}

	//Get code_challenge, and if client_id settings require use of PKCE, return an error if not respected.
	codeChallenge := r.URL.Query().Get(constants.PARAM_CODE_CHALLENGE)
	if codeChallenge == "" && authRequest.ClientId.ForceUseOfPKCE {
		return oauth2_errors.MissingCodeChallenge
	}

	codeChallengeMethod := CodeChallengeMethod(r.URL.Query().Get(constants.PARAM_CODE_CHALLENGE_METHOD))
	//If code_challenge_method is specified, then the value must be plain or S256
	if codeChallengeMethod != "" && codeChallengeMethod != CODE_CHALLENGE_METHOD_PLAIN && codeChallengeMethod != CODE_CHALLENGE_METHOD_S256 {
		return oauth2_errors.InvalidCodeChallenge
	}

	//If the code_challenge_method is not specified, but there's a code_challenge informed, so we use plain as default
	//For more details, see : https://tools.ietf.org/html/rfc7636#section-4.3
	if codeChallenge != "" && codeChallengeMethod == "" {
		codeChallengeMethod = CODE_CHALLENGE_METHOD_PLAIN
	}

	authRequest.codeChallenge = codeChallenge
	authRequest.codeChallengeMethod = codeChallengeMethod

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
	*/
	return nil
}

func handleImplicitFlowRequest(w http.ResponseWriter, r *http.Request, authRequest *models.AuthorizationRequest) error {

	//Initialize redirect_uri (optional query parameter)
	/*if redirectUri, err := initRedirectUri(r, authRequest.ClientId.AllowedRedirectUri, true); err != nil {
		return err
	} else {
		authRequest.RedirectUri = redirectUri
	}

	//TODO generate true jwt
	accessToken := "yoloooo"
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

	http.Redirect(w, r, location, http.StatusFound)*/

	return nil
}

/*
 *  On the authorization code flow, the redirect_uri is required : https://tools.ietf.org/html/rfc6749#section-4.1.1
 *  But on implicit flow, it is not mandatory as specified here : https://tools.ietf.org/html/rfc6749#section-4.2.1
 *  In such case we must ensure that the request come's from an allowed client uri https://tools.ietf.org/html/rfc6749#section-3.1.2
 */
/*func initRedirectUri(r *http.Request, allowedRedirectUris []string, isImplicit bool) (string, *oauth2_errors.Error) {

	//TODO : check wildcard uri (be as restrictive as possible)

	// If redirect_uri is not informed and current request is oauth2 implicit flow, then we get it from the settings.
	if redirectUri := r.URL.Query().Get(constants.PARAM_REDIRECT_URI); redirectUri == "" && isImplicit && len(allowedRedirectUris) == 1 {
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
	if isImplicit {
		return "", ImplicitFlowRedirectUriError
	} else {
		return "", RedirectUriError
	}
}

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
}*/
