package client

import (
	"oauth2-provider/models"
	"errors"
)

//Should be implemented later...
func FindAndLoadClientSettings(clientId string) (*models.ClientId, error) {

	//TODO add authorized oauth2 flows...

	//to implement later
	/*switch clientId {
	case "toto":
		return &ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback", "http://callback2"}, ForceUseOfPKCE: true}, nil
	case "titi":
		return &ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback", "http://callback2"}, ForceUseOfPKCE: false}, nil
	case "tutu":
		return &ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback"}}, nil
	default:
		return nil, &oauth2_errors.BadRequest{
			Reason:           oauth2_errors.ERROR_INVALID_CLIENT,
			ErrorDescription: oauth2_errors.DESC_INVALID_CLIENT,
			ErrorUri:         "https://tools.ietf.org/html/rfc6749#section-2.2",
		}
	}*/

	return nil, errors.New("Not implemented")
}
