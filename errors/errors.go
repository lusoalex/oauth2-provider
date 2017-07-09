package errors

import (
	"oauth2-provider/models"
)

const (
	ERROR_INVALID_REQUEST           = "invalid_request"
	ERROR_INVALID_CLIENT            = "invalid_client"
	ERROR_INVALID_GRANT             = "invalid_grant"
	ERROR_UNAUTHORIZED_CLIENT       = "unauthorized_client"
	ERROR_UNSUPPORTED_GRANT_TYPE    = "unsupported_grant_type"
	ERROR_INVALID_SCOPE             = "invalid_scope"
	ERROR_UNSUPPORTED_RESPONSE_TYPE = "unsupported_response_type"

	DESC_UNSUPPORTED_RESPONSE_TYPE = "Missing, unsupported or malformed required response_type parameter."
	DESC_UNSUPPORTED_GRANT_TYPE    = "Missing, unsupported or malformed required grant_type parameter."
	DESC_MISSING_CODE_CHALLENGE    = "Missing required code_challenger parameter."
	DESC_INVALID_CODE_CHALLENGE    = "Invalid code_challange_method parameter"
	DESC_INVALID_CLIENT            = "Missing or Unknown required client_id parameter."
	DESC_INVALID_REDIRECT_URI      = "Missing, invalid, or mismatching redirect_uri parameter."
)

var InvalidTypeError = &models.Error{
	Reason:           ERROR_UNSUPPORTED_RESPONSE_TYPE,
	ErrorDescription: DESC_UNSUPPORTED_RESPONSE_TYPE,
	ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-3.1.1",
}

var InvalidRedirectUri = &models.Error{
	Reason:           ERROR_INVALID_REQUEST,
	ErrorDescription: DESC_INVALID_REDIRECT_URI,
	ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-3.1.2",
}

var MissingCodeChallenge = &models.Error{
	Reason:           ERROR_INVALID_REQUEST,
	ErrorDescription: DESC_MISSING_CODE_CHALLENGE,
	ErrorUri:         "https://tools.ietf.org/html/rfc7636#section-4.4.1",
}

var InvalidCodeChallenge = &models.Error{
	Reason:           ERROR_INVALID_REQUEST,
	ErrorDescription: DESC_INVALID_CODE_CHALLENGE,
	ErrorUri:         "https://tools.ietf.org/html/rfc7636#section-4.3",
}

var UnsupportedGrantType = &models.Error{
	Reason:           ERROR_UNSUPPORTED_GRANT_TYPE,
	ErrorDescription: DESC_UNSUPPORTED_GRANT_TYPE,
	ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-5.2",
}

var MissingOrMalformedCodeVerifierFormat = &models.Error{
	Reason:           ERROR_INVALID_REQUEST,
	ErrorDescription: "Missing or malformed code_verifier parameter",
	ErrorUri:         "https://tools.ietf.org/html/rfc7636#section-4.1",
}

var InvalidCodeVerifier = &models.Error{
	Reason:           ERROR_INVALID_GRANT,
	ErrorDescription: "Invalid code_verifier parameter",
	ErrorUri:         "https://tools.ietf.org/html/rfc7636#section-4.6",
}

/*func (error *Oauth2Error) handle(w http.ResponseWriter, status int) {
	if errorMessage, err := json.Marshal(error); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(status)
		w.Write([]byte(errorMessage))
	}
}

func (error *models.Error) Handle(w http.ResponseWriter) {
	error.handle(w, http.Statusmodels.Error)
}

func (error *ForbiddenRequest) Handle(w http.ResponseWriter) {
	error.handle(w, http.StatusForbidden)
}*/

/*
func Oauth2NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage, _ := json.Marshal(&Oauth2models.Error{
		Reason: http.StatusText(http.StatusNotFound),
	})
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(errorMessage))
}

func Oauth2MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage, _ := json.Marshal(&Oauth2models.Error{
		Reason: http.StatusText(http.StatusMethodNotAllowed),
	})
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(errorMessage))
}
*/
