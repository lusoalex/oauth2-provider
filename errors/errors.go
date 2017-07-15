package errors

import (
	"oauth2-provider/models"
)

func UnauthorizedClient(grantType models.GrantType) *models.BadRequest {
	return BadRequest(
		"unauthorized_client",
		string(grantType)+" grant not allowed for this client",
		"https://tools.ietf.org/html/rfc6749#section-5.2",
	)
}

func InvalidRequest(desc, uri string) *models.BadRequest {
	return BadRequest("invalid_request", desc, uri)
}

func InvalidGrant(desc string) *models.BadRequest {
	return BadRequest("invalid_grant", desc, "https://tools.ietf.org/html/rfc6749#section-5.2")
}

func BadRequest(error, desc, uri string) *models.BadRequest {
	return &models.BadRequest{Oauth2Error: &models.Oauth2Error{
		Reason:           error,
		ErrorDescription: desc,
		ErrorUri:         uri,
	}}
}
