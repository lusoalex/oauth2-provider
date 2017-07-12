package errors

import (
	"oauth2-provider/models"
)

const (
	ERROR_INVALID_REQUEST           = "invalid_request"
	ERROR_INVALID_GRANT             = "invalid_grant"
	ERROR_UNAUTHORIZED_CLIENT       = "unauthorized_client"
	ERROR_UNSUPPORTED_GRANT_TYPE    = "unsupported_grant_type"
	ERROR_INVALID_SCOPE             = "invalid_scope"
	ERROR_UNSUPPORTED_RESPONSE_TYPE = "unsupported_response_type"

	DESC_UNSUPPORTED_RESPONSE_TYPE = "Missing, unsupported or malformed required response_type parameter."
	DESC_UNSUPPORTED_GRANT_TYPE    = "Missing, unsupported or malformed required grant_type parameter."
	DESC_MISSING_CODE_CHALLENGE    = "Missing required code_challenger parameter."
	DESC_INVALID_CODE_CHALLENGE    = "Invalid code_challange_method parameter"
	DESC_INVALID_REDIRECT_URI      = "Missing, invalid, or mismatching redirect_uri parameter."
)

var InvalidTypeError = &models.Oauth2Error{
	Reason:           ERROR_UNSUPPORTED_RESPONSE_TYPE,
	ErrorDescription: DESC_UNSUPPORTED_RESPONSE_TYPE,
	ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-3.1.1",
}

var InvalidRedirectUri = &models.Oauth2Error{
	Reason:           ERROR_INVALID_REQUEST,
	ErrorDescription: DESC_INVALID_REDIRECT_URI,
	ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-3.1.2",
}

var MissingCodeChallenge = &models.Oauth2Error{
	Reason:           ERROR_INVALID_REQUEST,
	ErrorDescription: DESC_MISSING_CODE_CHALLENGE,
	ErrorUri:         "https://tools.ietf.org/html/rfc7636#section-4.4.1",
}

var InvalidCodeChallenge = &models.Oauth2Error{
	Reason:           ERROR_INVALID_REQUEST,
	ErrorDescription: DESC_INVALID_CODE_CHALLENGE,
	ErrorUri:         "https://tools.ietf.org/html/rfc7636#section-4.3",
}

var UnsupportedGrantType = &models.Oauth2Error{
	Reason:           ERROR_UNSUPPORTED_GRANT_TYPE,
	ErrorDescription: DESC_UNSUPPORTED_GRANT_TYPE,
	ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-5.2",
}

var MissingOrMalformedCodeVerifierFormat = &models.Oauth2Error{
	Reason:           ERROR_INVALID_REQUEST,
	ErrorDescription: "Missing or malformed code_verifier parameter",
	ErrorUri:         "https://tools.ietf.org/html/rfc7636#section-4.1",
}

var InvalidCodeVerifier = &models.Oauth2Error{
	Reason:           ERROR_INVALID_GRANT,
	ErrorDescription: "Invalid code_verifier parameter",
	ErrorUri:         "https://tools.ietf.org/html/rfc7636#section-4.6",
}

func UnauthorizedClient(grantType models.GrantType) *models.Oauth2Error {
	return &models.Oauth2Error{
		Reason:           ERROR_UNAUTHORIZED_CLIENT,
		ErrorDescription: string(grantType) + " grant not allowed for this client",
		ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-5.2",
	}
}
