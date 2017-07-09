package client

import (
	"oauth2-provider/errors"
)

type ClientId struct {
	ClientId           string
	AllowedRedirectUri []string
	AllowedOauth2Flow  []Oauth2Flow
	ForceUseOfPKCE     bool
}

//Should be implemented later...
func FindAndLoadClientSettings(clientId string) (*ClientId, *errors.Oauth2Error) {

	//TODO add authorized oauth2 flows...

	//to implement later
	switch clientId {
	case "toto":
		return &ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback", "http://callback2"}, ForceUseOfPKCE: true}, nil
	case "titi":
		return &ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback", "http://callback2"}, ForceUseOfPKCE: false}, nil
	case "tutu":
		return &ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback"}}, nil
	default:
		return nil, errors.ClientIdError
	}
}
