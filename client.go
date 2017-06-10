package oauth2Provider

type ClientId struct {
	ClientId           string
	AllowedRedirectUri []string
	AllowedOauth2Flow  []Oauth2Flow
	ForceUseOfPKCE     bool
}

//Should be implemented later...
func findAndLoadClientSettings(clientId string) (*ClientId, *Oauth2Error) {

	//to implement later
	if clientId == "toto" {
		return &ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback"}, ForceUseOfPKCE: true}, nil
	} else if clientId == "titi" {
		return &ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback"}, ForceUseOfPKCE: false}, nil
	} else {
		return nil, NewClientIdError()
	}
}
