package oauth2Provider

type ClientId struct {
	ClientId           string
	AllowedRedirectUri []string
	AllowedOauth2Flow  []Oauth2Flow
	ForceUseOfPKCE     bool
}

//Should be implemented later...
func findAndLoadClientSettings(clientId string) (*ClientId, *Oauth2Error) {

	//TODO add authorized oauth2 flows...

	//to implement later
	if clientId == "toto" {
		return &ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback", "http://callback2"}, ForceUseOfPKCE: true}, nil
	} else if clientId == "titi" {
		return &ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback", "http://callback2"}, ForceUseOfPKCE: false}, nil
	} else if clientId == "tutu" {
		return &ClientId{ClientId: clientId, AllowedRedirectUri: []string{"http://callback"}}, nil
	} else {
		return nil, NewClientIdError()
	}
}
