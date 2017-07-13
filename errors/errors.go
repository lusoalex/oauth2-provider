package errors

import (
	"oauth2-provider/models"
)

const (
	ERROR_INVALID_REQUEST     = "invalid_request"
	ERROR_UNAUTHORIZED_CLIENT = "unauthorized_client"
)

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
