package client

import (
	"errors"
	"oauth2-provider/models"
)

const (
	ERROR_INVALID_CLIENT = "invalid_client"
	DESC_INVALID_CLIENT  = "Missing or Unknown required client_id parameter."
)

var INVALID_CLIENT_ID = &models.Error{
	Reason:           ERROR_INVALID_CLIENT,
	ErrorDescription: DESC_INVALID_CLIENT,
	ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-2.2",
}

//TODO Should be implemented later...
func GetClientInformations(clientId string) (*models.ClientId, error) {

	switch clientId {
	case "toto":
		return &models.ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback", "http://callback2"}, AllowedGrantType: []models.GrantType{"implicit", "authorization_code"}, ForceUseOfPKCE: true}, nil
	case "titi":
		return &models.ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback", "http://callback2"}, AllowedGrantType: []models.GrantType{"implicit", "authorization_code"}, ForceUseOfPKCE: false}, nil
	case "tutu":
		return &models.ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback"}, AllowedGrantType: []models.GrantType{"implicit", "authorization_code"}}, nil
	default:
		return nil, errors.New(ERROR_INVALID_CLIENT)
	}
}
