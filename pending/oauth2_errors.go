package oauth2Provider

import (
	"encoding/json"
	"net/http"
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

/*
 * Using fields as specified in https://tools.ietf.org/html/rfc6749#section-5.2
 */
type Oauth2Error struct {
	status           int    `json:"omit"`                        //to omit in the body response
	Reason           string `json:"error"`                       //required
	ErrorDescription string `json:"error_description,omitempty"` //Optional
	ErrorUri         string `json:"error_uri,omitempty"`         //Optional
	State            string `json:"state,omitempty"`             //Required if present into the request.
}

var ResponseTypeError = &Oauth2Error{
	status:           http.StatusBadRequest,
	Reason:           ERROR_UNSUPPORTED_RESPONSE_TYPE,
	ErrorDescription: DESC_UNSUPPORTED_RESPONSE_TYPE,
	ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-3.1.1",
}

var ClientIdError = &Oauth2Error{
	status:           http.StatusBadRequest,
	Reason:           ERROR_INVALID_CLIENT,
	ErrorDescription: DESC_INVALID_CLIENT,
	ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-2.2",
}

var ImplicitFlowRedirectUriError = &Oauth2Error{
	status:           http.StatusBadRequest,
	Reason:           ERROR_INVALID_REQUEST,
	ErrorDescription: DESC_INVALID_REDIRECT_URI,
	ErrorUri: "https://tools.ietf.org/html/rfc6749#section-4.2.2.1",
}

var RedirectUriError = &Oauth2Error{
	status:           http.StatusBadRequest,
	Reason:           ERROR_INVALID_REQUEST,
	ErrorDescription: DESC_INVALID_REDIRECT_URI,
	ErrorUri: "https://tools.ietf.org/html/rfc6749#section-4.1.2.1",
}

var CodeChallengeError = &Oauth2Error{
	status:           http.StatusBadRequest,
	Reason:           ERROR_INVALID_REQUEST,
	ErrorDescription: DESC_MISSING_CODE_CHALLENGE,
	ErrorUri:         "https://tools.ietf.org/html/rfc7636#section-4.4.1",
}

var CodeChallengeMethodError = &Oauth2Error{
	status:           http.StatusBadRequest,
	Reason:           ERROR_INVALID_REQUEST,
	ErrorDescription: DESC_INVALID_CODE_CHALLENGE,
	ErrorUri:         "https://tools.ietf.org/html/rfc7636#section-4.3",
}

var GrantTypeError = &Oauth2Error{
	status:           http.StatusBadRequest,
	Reason:           ERROR_UNSUPPORTED_GRANT_TYPE,
	ErrorDescription: DESC_UNSUPPORTED_GRANT_TYPE,
	ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-5.2",
}

var CodeVerifierFormatError = &Oauth2Error{
	status:           http.StatusBadRequest,
	Reason:           ERROR_INVALID_REQUEST,
	ErrorDescription: "Missing or malformed code_verifier parameter",
	ErrorUri:         "https://tools.ietf.org/html/rfc7636#section-4.1",
}

var CodeVerifierError = &Oauth2Error{
	status:           http.StatusBadRequest,
	Reason:           ERROR_INVALID_GRANT,
	ErrorDescription: "Invalid code_verifier parameter",
	ErrorUri:         "https://tools.ietf.org/html/rfc7636#section-4.6",
}

func (error *Oauth2Error) Handle(w http.ResponseWriter) {
	if errorMessage, err := json.Marshal(error) ; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(error.status)
		w.Write([]byte(errorMessage))
	}
}

func Oauth2NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage, _ := json.Marshal(&Oauth2Error{
		Reason: http.StatusText(http.StatusNotFound),
	})
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(errorMessage))
}

func Oauth2MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage, _ := json.Marshal(&Oauth2Error{
		Reason: http.StatusText(http.StatusMethodNotAllowed),
	})
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(errorMessage))
}
