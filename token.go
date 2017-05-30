package oauth2Provider

type Token struct {
	AccessToken *string `json:"access_token"`
	RefreshToken *string `json:"refresh_token,omitempty"`
}