package models

type ResponseType string
type CodeChallengeMethod string

const (
	RESPONSE_TYPE_CODE  ResponseType = "code"
	RESPONSE_TYPE_TOKEN ResponseType = "token"

	CODE_CHALLENGE_METHOD_PLAIN CodeChallengeMethod = "plain"
	CODE_CHALLENGE_METHOD_S256  CodeChallengeMethod = "S256"
)

type AuthorizationRequest struct {
	ClientId            ClientId
	ResponseType        ResponseType
	RedirectUri         string
	Scope               string
	State               string
	codeChallenge       string
	codeChallengeMethod CodeChallengeMethod
	Code                string
}

