package client

import (
	"errors"
	errors2 "oauth2-provider/errors"
	"oauth2-provider/models"
)

//Should be implemented later...
func GetClientInformations(clientId string) (*models.ClientId, error) {

	//TODO remove mock and implement

	switch clientId {
	case "toto":
		return &models.ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback", "http://callback2"}, AllowedGrantType: []models.GrantType{"implicit", "authorization_code"}, ForceUseOfPKCE: true}, nil
	case "titi":
		return &models.ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback", "http://callback2"}, AllowedGrantType: []models.GrantType{"implicit", "authorization_code"}, ForceUseOfPKCE: false}, nil
	case "tutu":
		return &models.ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback"}, AllowedGrantType: []models.GrantType{"implicit", "authorization_code"}}, nil
	default:
		return nil, errors.New(errors2.ERROR_INVALID_CLIENT)
	}
}
