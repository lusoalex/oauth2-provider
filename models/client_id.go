package models

type ClientId struct {
	ClientId           string
	AllowedRedirectUri []string
	AllowedGrantType   []GrantType
	ForceUseOfPKCE     bool
}
