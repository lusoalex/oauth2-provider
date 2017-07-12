package client

import (
	"oauth2-provider/models"
)

const (
	INVALID_CLIENT_ERROR = "invalid_client"
	INVALID_CLIENT_DESC  = "Missing or Unknown required client_id parameter."
	INVALID_CLIENT_URI   = "https://tools.ietf.org/html/rfc6749#section-2.2"
)

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
		return nil, &models.BadRequest{
			Oauth2Error: &models.Oauth2Error{
				Reason:           INVALID_CLIENT_ERROR,
				ErrorDescription: INVALID_CLIENT_DESC,
				ErrorUri:         INVALID_CLIENT_URI,
			}}
	}
}
