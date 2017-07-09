package constants

const (
	CONTENT_TYPE      = "Content-Type"
	CONTENT_TYPE_JSON = "application/json"
	CONTENT_TYPE_HTML = "text/html"
	CONTENT_TYPE_TEXT = "text/plain"

	//Request query parameters needed of oauth2 flows https://tools.ietf.org/html/rfc6749#appendix-A
	PARAM_CLIENT_ID     = "client_id"
	PARAM_CLIENT_SECRET = "client_secret"
	PARAM_RESPONSE_TYPE = "response_type"
	PARAM_REDIRECT_URI  = "redirect_uri"
	PARAM_SCOPE         = "scope"
	PARAM_STATE         = "state"
	PARAM_GRANT_TYPE    = "grant_type"
	PARAM_CODE          = "code"
	PARAM_USERNAME      = "username"
	PARAM_PASSWORD      = "password"
	PARAM_ACCESS_TOKEN  = "access_token"
	PARAM_TOKEN_TYPE    = "token_type"

	//Additional query parameter as specified here : https://tools.ietf.org/html/rfc7636#section-6.1
	PARAM_CODE_CHALLENGE        = "code_challenge"
	PARAM_CODE_CHALLENGE_METHOD = "code_challenge_method"
	PARAM_CODE_VERIFIER         = "code_verifier"
)
