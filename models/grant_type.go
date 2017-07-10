package models

type GrantType string

const (
	GRANT_TYPE_IMPLICIT                             GrantType = "implicit"
	GRANT_TYPE_AUTHORIZATION_CODE                   GrantType = "authorization_code"
	GRANT_TYPE_REFRESH_TOKEN                        GrantType = "refresh_token"
	GRANT_TYPE_CLIENT_CREDENTIALS                   GrantType = "client_credentials"
	GRANT_TYPE_RESSOURCE_OWNER_PASSWORD_CREDENTIALS GrantType = "password"
	GRANT_TYPE_EXTENSION_GRANT                      GrantType = "extension_grant" //how to implement? https://tools.ietf.org/html/rfc6749#section-4.5
)
