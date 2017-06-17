package oauth2Provider

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type AuthorizeTestCase struct {
	ResponseType, ClientId, RedirectUri, CodeChallenge, CodeChallengeMethod string
}

const (
	CLIENT_ID_WITH_PKCE    = "toto"
	CLIENT_ID_WITHOUT_PKCE = "titi"
	VALID_REDIRECT_URI     = "http://callback"
	PLAIN_CODE_CHALLENGE   = "codeChallenge"
	S256_CODE_CHALLENGE    = "wxcqdsqd"
)

var CodeValidWithDefaultPlainPKCETestCase = &AuthorizeTestCase{ResponseType: string(RESPONSE_TYPE_CODE), ClientId: CLIENT_ID_WITH_PKCE, RedirectUri: VALID_REDIRECT_URI, CodeChallenge: PLAIN_CODE_CHALLENGE}
var CodeValidWithPlainPKCETestCase = &AuthorizeTestCase{ResponseType: string(RESPONSE_TYPE_CODE), ClientId: CLIENT_ID_WITH_PKCE, RedirectUri: VALID_REDIRECT_URI, CodeChallengeMethod: string(CODE_CHALLENGE_METHOD_PLAIN), CodeChallenge: PLAIN_CODE_CHALLENGE}
var CodeValidWithS256PKCETestCase = &AuthorizeTestCase{ResponseType: string(RESPONSE_TYPE_CODE), ClientId: CLIENT_ID_WITH_PKCE, RedirectUri: VALID_REDIRECT_URI, CodeChallengeMethod: string(CODE_CHALLENGE_METHOD_S256), CodeChallenge: S256_CODE_CHALLENGE}
var CodeValidWithoutPKCETestCase = &AuthorizeTestCase{ResponseType: string(RESPONSE_TYPE_CODE), ClientId: CLIENT_ID_WITHOUT_PKCE, RedirectUri: VALID_REDIRECT_URI}
var CodeValidWithOptionalPKCETestCase = &AuthorizeTestCase{ResponseType: string(RESPONSE_TYPE_CODE), ClientId: CLIENT_ID_WITHOUT_PKCE, RedirectUri: VALID_REDIRECT_URI, CodeChallenge: PLAIN_CODE_CHALLENGE}
var CodeMissingResponseTypeTestCase = &AuthorizeTestCase{ClientId: CLIENT_ID_WITH_PKCE, RedirectUri: VALID_REDIRECT_URI}
var CodeMissingClientIdTestCase = &AuthorizeTestCase{ResponseType: string(RESPONSE_TYPE_CODE), RedirectUri: VALID_REDIRECT_URI}
var CodeMissingRedirectUriTestCase = &AuthorizeTestCase{ResponseType: string(RESPONSE_TYPE_CODE), ClientId: CLIENT_ID_WITH_PKCE}
var CodeMissingCodeChallengeTestCse = &AuthorizeTestCase{ResponseType: string(RESPONSE_TYPE_CODE), ClientId: CLIENT_ID_WITH_PKCE, RedirectUri: VALID_REDIRECT_URI}
var CodeInvalidCodeChallengeMethodTestCase = &AuthorizeTestCase{ResponseType: string(RESPONSE_TYPE_CODE), ClientId: CLIENT_ID_WITH_PKCE, RedirectUri: VALID_REDIRECT_URI, CodeChallengeMethod: "fail", CodeChallenge: S256_CODE_CHALLENGE}
var CodeInvalidResponseTypeTestCase = &AuthorizeTestCase{ResponseType: "blabla", ClientId: CLIENT_ID_WITH_PKCE, RedirectUri: VALID_REDIRECT_URI}
var CodeInvalidClientIdTestCase = &AuthorizeTestCase{ResponseType: string(RESPONSE_TYPE_CODE), ClientId: "badClientName", RedirectUri: VALID_REDIRECT_URI}
var CodeInvalidRedirectUriTestCase = &AuthorizeTestCase{ResponseType: string(RESPONSE_TYPE_CODE), ClientId: CLIENT_ID_WITH_PKCE, RedirectUri: "http://fail/back"}

var ImplicitValidTestCase = &AuthorizeTestCase{ResponseType: string(RESPONSE_TYPE_TOKEN), ClientId: CLIENT_ID_WITH_PKCE, RedirectUri: VALID_REDIRECT_URI}
var ImplicitWithoutRedirectUriTestCase = &AuthorizeTestCase{ResponseType: string(RESPONSE_TYPE_TOKEN), ClientId: CLIENT_ID_WITH_PKCE}
var ImpicitErrorTestCase = &AuthorizeTestCase{ResponseType: string(RESPONSE_TYPE_TOKEN), ClientId: CLIENT_ID_WITH_PKCE, RedirectUri: "http://fail/back"}

func TestHandleAuthorizationRequest(t *testing.T) {

	testCases := []DataTestCase{
		{name: "Valid authorization code with PKCE and default challenge_method flow", url: buildTestCaseUrl(CodeValidWithDefaultPlainPKCETestCase), status: http.StatusOK},
		{name: "Valid authorization code with PKCE and plain challenge_method flow", url: buildTestCaseUrl(CodeValidWithPlainPKCETestCase), status: http.StatusOK},
		{name: "Valid authorization code with PKCE and S256 challenge_method flow", url: buildTestCaseUrl(CodeValidWithS256PKCETestCase), status: http.StatusOK},
		{name: "Valid authorization code without PKCE flow", url: buildTestCaseUrl(CodeValidWithoutPKCETestCase), status: http.StatusOK},
		{name: "Valid authorization code without PKCE but given anyway", url: buildTestCaseUrl(CodeValidWithOptionalPKCETestCase), status: http.StatusOK},
		{name: "Valid implicit flow", url: buildTestCaseUrl(ImplicitValidTestCase), status: http.StatusOK},
		{name: "Valid implicit flow without redirect uri", url: buildTestCaseUrl(ImplicitWithoutRedirectUriTestCase), status: http.StatusOK},
		{name: "Missing response_type", url: buildTestCaseUrl(CodeMissingResponseTypeTestCase), status: http.StatusBadRequest, error: DESC_UNSUPPORTED_RESPONSE_TYPE},
		{name: "Missing client_id", url: buildTestCaseUrl(CodeMissingClientIdTestCase), status: http.StatusBadRequest, error: DESC_INVALID_CLIENT},
		{name: "Missing code redirect_uri", url: buildTestCaseUrl(CodeMissingRedirectUriTestCase), status: http.StatusBadRequest, error: DESC_INVALID_REDIRECT_URI},
		{name: "Missing code code_challenge", url: buildTestCaseUrl(CodeMissingCodeChallengeTestCse), status: http.StatusBadRequest, error: DESC_MISSING_CODE_CHALLENGE},
		{name: "Invalid code challenge method", url: buildTestCaseUrl(CodeInvalidCodeChallengeMethodTestCase), status: http.StatusBadRequest, error: DESC_INVALID_CODE_CHALLENGE},
		{name: "Invalid response_type", url: buildTestCaseUrl(CodeInvalidResponseTypeTestCase), status: http.StatusBadRequest, error: DESC_UNSUPPORTED_RESPONSE_TYPE},
		{name: "Invalid client_id", url: buildTestCaseUrl(CodeInvalidClientIdTestCase), status: http.StatusBadRequest, error: DESC_INVALID_CLIENT},
		{name: "Invalid redirect_uri", url: buildTestCaseUrl(CodeInvalidRedirectUriTestCase), status: http.StatusBadRequest, error: DESC_INVALID_REDIRECT_URI},
		{name: "Invalid implicit flow", url: buildTestCaseUrl(ImpicitErrorTestCase), status: http.StatusBadRequest, error: DESC_INVALID_REDIRECT_URI},
	}

	//And we set it the handler we aim to test
	handler := http.HandlerFunc(AuthorizationRequestHandler)

	for _, test := range testCases {

		//Initialize request
		req, err := http.NewRequest(http.MethodGet, test.url, nil)

		if err != nil {
			t.Fatal(err)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		//Check the status code is what we expect.
		assertStatus(t, &test, rr.Code)

		//If status is not ok, and error message is not as expected
		if rr.Code != http.StatusOK {
			assertBodyErrorMessage(t, &test, rr.Body.String())
		}
	}
}

/* SOME utils function to build test case values*/
func buildTestCaseUrl(testCase *AuthorizeTestCase) string {

	uri, _ := url.Parse("/authorize")
	query := uri.Query()
	query.Add(PARAM_CLIENT_ID, testCase.ClientId)
	query.Add(PARAM_RESPONSE_TYPE, testCase.ResponseType)
	query.Add(PARAM_REDIRECT_URI, testCase.RedirectUri)
	query.Add(PARAM_CODE_CHALLENGE, testCase.CodeChallenge)
	query.Add(PARAM_CODE_CHALLENGE_METHOD, testCase.CodeChallengeMethod)

	uri.RawQuery = query.Encode()
	return uri.String()
}
