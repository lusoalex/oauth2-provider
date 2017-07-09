package models

type ClientId struct {
	ClientId           string
	AllowedRedirectUri []string
	//AllowedOauth2Flow  []handlers.Oauth2Flow
	ForceUseOfPKCE     bool
}
