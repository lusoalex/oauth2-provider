package oauth2Provider

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type DataTestCase struct {
	name, url string
	status    int
	error     string
}

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
		{"valid authorization code with PKCE and default challenge_method flow", buildTestCaseUrl(CodeValidWithDefaultPlainPKCETestCase), http.StatusOK, ""},
		{"valid authorization code with PKCE and plain challenge_method flow", buildTestCaseUrl(CodeValidWithPlainPKCETestCase), http.StatusOK, ""},
		{"valid authorization code with PKCE and S256 challenge_method flow", buildTestCaseUrl(CodeValidWithS256PKCETestCase), http.StatusOK, ""},
		{"valid authorization code without PKCE flow", buildTestCaseUrl(CodeValidWithoutPKCETestCase), http.StatusOK, ""},
		{"valid implicit flow", buildTestCaseUrl(ImplicitValidTestCase), http.StatusOK, ""},
		{"valid implicit flow without redirect uri", buildTestCaseUrl(ImplicitWithoutRedirectUriTestCase), http.StatusOK, ""},
		{"missing response_type", buildTestCaseUrl(CodeMissingResponseTypeTestCase), http.StatusBadRequest, DESC_UNSUPPORTED_RESPONSE_TYPE},
		{"missing client_id", buildTestCaseUrl(CodeMissingClientIdTestCase), http.StatusBadRequest, DESC_INVALID_CLIENT},
		{"missing code redirect_uri", buildTestCaseUrl(CodeMissingRedirectUriTestCase), http.StatusBadRequest, DESC_INVALID_REDIRECT_URI},
		{"missing code code_challenge", buildTestCaseUrl(CodeMissingCodeChallengeTestCse), http.StatusBadRequest, DESC_MISSING_CODE_CHALLENGE},
		{"invalid code challenge method", buildTestCaseUrl(CodeInvalidCodeChallengeMethodTestCase), http.StatusBadRequest, DESC_INVALID_CODE_CHALLENGE},
		{"invalid response_type", buildTestCaseUrl(CodeInvalidResponseTypeTestCase), http.StatusBadRequest, DESC_UNSUPPORTED_RESPONSE_TYPE},
		{"invalid client_id", buildTestCaseUrl(CodeInvalidClientIdTestCase), http.StatusBadRequest, DESC_INVALID_CLIENT},
		{"invalid redirect_uri", buildTestCaseUrl(CodeInvalidRedirectUriTestCase), http.StatusBadRequest, DESC_INVALID_REDIRECT_URI},
		{"invalid implicit flow", buildTestCaseUrl(ImpicitErrorTestCase), http.StatusBadRequest, DESC_INVALID_REDIRECT_URI},
	}

	//And we set it the handler we aim to test
	handler := http.HandlerFunc(handleAuthorizationRequest)

	for _, test := range testCases {

		//Initialize request
		req, err := http.NewRequest("GET", test.url, nil)

		if err != nil {
			t.Fatal(err)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		//Status response
		status := rr.Code

		//Check the status code is what we expect.
		if status != test.status {
			t.Errorf("Test[%v] : handler returned wrong status code: got %v want %v",
				test.name, status, test.status)
		}

		//If status is not ok, and error message is not as expected
		if body := rr.Body.String(); status != http.StatusOK && !strings.Contains(body, test.error) {
			t.Errorf("Test[%v] : handler body response do not contains expected error message [%v], got [%v]",
				test.name, test.error, body)
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
