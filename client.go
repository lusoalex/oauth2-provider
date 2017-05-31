package oauth2Provider

import "errors"

type ClientId struct {
	ClientId           string
	AllowedRedirectUri []string
	AllowedOauth2Flow  []Oauth2Flow
}

const (
	INVALID_CLIENT_ID = "Unknown or invalid client_id"
)

//Should be implemented later...
func findAndLoadClientSettings(clientId string) (*ClientId, error) {

	//to implement later
	if clientId == "toto" {
		return &ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback"}}, nil
	} else {
		return nil, errors.New(INVALID_CLIENT_ID)
	}
}
