package oauth2Provider

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const (
	VALID_PATH          = "/authorize"
	VALID_CLIENT_ID     = "toto"
	VALID_RESPONSE_TYPE = "code"
	VALID_REDIRECT_URI  = "http://callback"
)

type DataTestCase struct {
	name, url string
	status    int
	error     string
}

type TestCase struct {
	Path, ResponseType, ClientId, RedirectUri string
}

var ValidAuthorizationCodeTestCase = &TestCase{Path: VALID_PATH, ResponseType: VALID_RESPONSE_TYPE, ClientId: VALID_CLIENT_ID, RedirectUri: VALID_REDIRECT_URI}
var MissingResponseTypeTestCase = &TestCase{Path: VALID_PATH, ClientId: VALID_CLIENT_ID, RedirectUri: VALID_REDIRECT_URI}
var MissingClientIdTestCase = &TestCase{Path: VALID_PATH, ResponseType: VALID_RESPONSE_TYPE, RedirectUri: VALID_REDIRECT_URI}
var MissingRedirectUriTestCase = &TestCase{Path: VALID_PATH, ResponseType: VALID_RESPONSE_TYPE, ClientId: VALID_CLIENT_ID}
var InvalidResponseTypeTestCase = &TestCase{Path: VALID_PATH, ResponseType: "blabla", ClientId: VALID_CLIENT_ID, RedirectUri: VALID_REDIRECT_URI}
var InvalidClientIdTestCase = &TestCase{Path: VALID_PATH, ResponseType: VALID_RESPONSE_TYPE, ClientId: "badClientName", RedirectUri: VALID_REDIRECT_URI}
var InvalidRedirectUriTestCase = &TestCase{Path: VALID_PATH, ResponseType: VALID_RESPONSE_TYPE, ClientId: VALID_CLIENT_ID, RedirectUri: "http://fail/back"}

func TestHandleAuthorizationRequest(t *testing.T) {

	testCases := []DataTestCase{
		{"valid path", buildTestCaseUrl(ValidAuthorizationCodeTestCase), http.StatusOK, ""},
		{"missing response_type", buildTestCaseUrl(MissingResponseTypeTestCase), http.StatusBadRequest, DESC_UNSUPPORTED_RESPONSE_TYPE},
		{"missing client_id", buildTestCaseUrl(MissingClientIdTestCase), http.StatusBadRequest, DESC_INVALID_CLIENT},
		{"missing redirect_uri", buildTestCaseUrl(MissingRedirectUriTestCase), http.StatusBadRequest, DESC_INVALID_REDIRECT_URI},
		{"invalid response_type", buildTestCaseUrl(InvalidResponseTypeTestCase), http.StatusBadRequest, DESC_UNSUPPORTED_RESPONSE_TYPE},
		{"invalid client_id", buildTestCaseUrl(InvalidClientIdTestCase), http.StatusBadRequest, DESC_INVALID_CLIENT},
		{"invalid redirect_uri", buildTestCaseUrl(InvalidRedirectUriTestCase), http.StatusBadRequest, DESC_INVALID_REDIRECT_URI},
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
func buildTestCaseUrl(testCase *TestCase) string {

	uri, _ := url.Parse(testCase.Path)
	query := uri.Query()
	query.Add(PARAM_CLIENT_ID, testCase.ClientId)
	query.Add(PARAM_RESPONSE_TYPE, testCase.ResponseType)
	query.Add(PARAM_REDIRECT_URI, testCase.RedirectUri)

	uri.RawQuery = query.Encode()
	return uri.String()
}
