package errors

import (
	"oauth2-provider/models"
)

const (
	ERROR_INVALID_REQUEST           = "invalid_request"
	ERROR_INVALID_GRANT             = "invalid_grant"
	ERROR_UNAUTHORIZED_CLIENT       = "unauthorized_client"
	ERROR_INVALID_SCOPE             = "invalid_scope"
	ERROR_UNSUPPORTED_RESPONSE_TYPE = "unsupported_response_type"

	DESC_UNSUPPORTED_RESPONSE_TYPE = "Missing, unsupported or malformed required response_type parameter."
)

var InvalidTypeError = &models.Oauth2Error{
	Reason:           ERROR_UNSUPPORTED_RESPONSE_TYPE,
	ErrorDescription: DESC_UNSUPPORTED_RESPONSE_TYPE,
	ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-3.1.1",
}

/*
var MissingOrMalformedCodeVerifierFormat = &models.Oauth2Error{
	Reason:           ERROR_INVALID_REQUEST,
	ErrorDescription: "Missing or malformed code_verifier parameter",
	ErrorUri:         "https://tools.ietf.org/html/rfc7636#section-4.1",
}
*/

var InvalidCodeVerifier = &models.Oauth2Error{
	Reason:           ERROR_INVALID_GRANT,
	ErrorDescription: "Invalid code_verifier parameter",
	ErrorUri:         "https://tools.ietf.org/html/rfc7636#section-4.6",
}

func UnauthorizedClient(grantType models.GrantType) *models.BadRequest {
	return &models.BadRequest{Oauth2Error: &models.Oauth2Error{
		Reason:           ERROR_UNAUTHORIZED_CLIENT,
		ErrorDescription: string(grantType) + " grant not allowed for this client",
		ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-5.2",
	}}
}

func InvalidRequest(desc, uri string) *models.BadRequest {
	return &models.BadRequest{Oauth2Error: &models.Oauth2Error{
		Reason:           ERROR_INVALID_REQUEST,
		ErrorDescription: desc,
		ErrorUri:         uri,
	}}
}
