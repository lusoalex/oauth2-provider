package oauth2Provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type AuthorizeTestCase struct {
	ResponseType, ClientId, RedirectUri, State, CodeChallenge, CodeChallengeMethod string
}

const (
	CLIENT_ID_CODE_WITH_PKCE    = "toto"
	CLIENT_ID_CODE_WITHOUT_PKCE = "titi"
	CLIENT_ID_IMPLICIT          = "tutu"
	VALID_REDIRECT_URI          = "http://callback"
	PLAIN_CODE_CHALLENGE        = "codeChallenge"
	S256_CODE_CHALLENGE         = "wxcqdsqd"
)

//Main handler we aim to test
var handler = http.HandlerFunc(AuthorizationRequestHandler)

/**************************************************/
/*                                                */
/*       AUTHORIZATION CODE FLOW TEST CASES       */
/*                                                */
/**************************************************/
func TestCodeFlowWithForcedPkceAndDefaultPlainChallengeMethod(t *testing.T) {
	//Build a valid authorize url for a client requiring a PKCE, but do not inform code_challenge method, default "plain" code_challenge should be used
	codeRequest := &AuthorizeTestCase{
		ResponseType:  string(RESPONSE_TYPE_CODE),
		ClientId:      CLIENT_ID_CODE_WITH_PKCE,
		RedirectUri:   VALID_REDIRECT_URI,
		CodeChallenge: PLAIN_CODE_CHALLENGE,
	}
	testValidAuthorizationCodeFlow(t, "Valid authorization code flow with forced PKCE and default challenge_method", codeRequest)
}

func TestCodeFlowWithForcedPkceAndPlainChallengeMethod(t *testing.T) {
	//Build a valid authorize url for a client requiring a PKCE giving a "plain" code_challenge method
	codeRequest := &AuthorizeTestCase{
		ResponseType:        string(RESPONSE_TYPE_CODE),
		ClientId:            CLIENT_ID_CODE_WITH_PKCE,
		RedirectUri:         VALID_REDIRECT_URI,
		CodeChallengeMethod: string(CODE_CHALLENGE_METHOD_PLAIN),
		CodeChallenge:       PLAIN_CODE_CHALLENGE,
	}
	testValidAuthorizationCodeFlow(t, "Valid authorization code flow with forced PKCE and plain challenge_method", codeRequest)
}

func TestCodeFlowWithForcedPkceAndS256ChallengeMethod(t *testing.T) {
	//Build a valid authorize url for a client requiring a PKCE giving a "S256" code_challenge method
	codeRequest := &AuthorizeTestCase{
		ResponseType:        string(RESPONSE_TYPE_CODE),
		ClientId:            CLIENT_ID_CODE_WITH_PKCE,
		RedirectUri:         VALID_REDIRECT_URI,
		CodeChallengeMethod: string(CODE_CHALLENGE_METHOD_S256),
		CodeChallenge:       S256_CODE_CHALLENGE,
	}
	testValidAuthorizationCodeFlow(t, "Valid authorization code flow with forced PKCE and S256 challenge_method", codeRequest)
}

func TestCodeFlowWithoutPkce(t *testing.T) {
	//Build a valid authorize url for a client not requiring a PKCE
	codeRequest := &AuthorizeTestCase{
		ResponseType: string(RESPONSE_TYPE_CODE),
		ClientId:     CLIENT_ID_CODE_WITHOUT_PKCE,
		RedirectUri:  VALID_REDIRECT_URI,
	}
	testValidAuthorizationCodeFlow(t, "Valid authorization code flow without PKCE", codeRequest)
}

func TestCodeFlowWithOptionalPkce(t *testing.T) {
	//Build a valid authorize url for a client not requiring a PKCE but giving it anyway
	codeRequest := &AuthorizeTestCase{
		ResponseType:  string(RESPONSE_TYPE_CODE),
		ClientId:      CLIENT_ID_CODE_WITHOUT_PKCE,
		RedirectUri:   VALID_REDIRECT_URI,
		CodeChallenge: PLAIN_CODE_CHALLENGE,
	}
	testValidAuthorizationCodeFlow(t, "Valid authorization code flow with optional PKCE provided", codeRequest)
}

func TestCodeWithMissingResponseType(t *testing.T) {
	codeRequest := &AuthorizeTestCase{
		ClientId:    CLIENT_ID_CODE_WITH_PKCE,
		RedirectUri: VALID_REDIRECT_URI,
	}
	testInvalidAuthorizationRequest(t, "Code request with missing response_type", codeRequest, DESC_UNSUPPORTED_RESPONSE_TYPE)
}

func TestCodeWithMissingClientId(t *testing.T) {
	codeRequest := &AuthorizeTestCase{
		ResponseType: string(RESPONSE_TYPE_CODE),
		RedirectUri:  VALID_REDIRECT_URI,
	}
	testInvalidAuthorizationRequest(t, "Code request with missing client_id", codeRequest, DESC_INVALID_CLIENT)
}

func TestCodeWithMissingRedirectUri(t *testing.T) {
	codeRequest := &AuthorizeTestCase{
		ResponseType: string(RESPONSE_TYPE_CODE),
		ClientId:     CLIENT_ID_CODE_WITH_PKCE,
	}
	testInvalidAuthorizationRequest(t, "Code request with missing redirect_uri", codeRequest, DESC_INVALID_REDIRECT_URI)
}

func TestCodeWithMissingCodeChallenge(t *testing.T) {
	codeRequest := &AuthorizeTestCase{
		ResponseType: string(RESPONSE_TYPE_CODE),
		ClientId:     CLIENT_ID_CODE_WITH_PKCE,
		RedirectUri:  VALID_REDIRECT_URI,
	}
	testInvalidAuthorizationRequest(t, "Code request with missing code_challenge", codeRequest, DESC_MISSING_CODE_CHALLENGE)
}

func TestCodeWithInvalidCodeChallengeMethod(t *testing.T) {
	codeRequest := &AuthorizeTestCase{
		ResponseType:        string(RESPONSE_TYPE_CODE),
		ClientId:            CLIENT_ID_CODE_WITH_PKCE,
		RedirectUri:         VALID_REDIRECT_URI,
		CodeChallengeMethod: "fail",
		CodeChallenge:       S256_CODE_CHALLENGE,
	}
	testInvalidAuthorizationRequest(t, "Code request with invalid code_challenge_method", codeRequest, DESC_INVALID_CODE_CHALLENGE)
}

func TestCodeWithInvalidResponseType(t *testing.T) {
	codeRequest := &AuthorizeTestCase{
		ResponseType: "blabla",
		ClientId:     CLIENT_ID_CODE_WITH_PKCE,
		RedirectUri:  VALID_REDIRECT_URI,
	}
	testInvalidAuthorizationRequest(t, "Code request with invalid response_type", codeRequest, DESC_UNSUPPORTED_RESPONSE_TYPE)
}

func TestCodeWithInvalidClientId(t *testing.T) {
	codeRequest := &AuthorizeTestCase{
		ResponseType: string(RESPONSE_TYPE_CODE),
		ClientId:     "badClientName",
		RedirectUri:  VALID_REDIRECT_URI,
	}
	testInvalidAuthorizationRequest(t, "Code request with invalid response_type", codeRequest, DESC_INVALID_CLIENT)
}

func TestCodeWithInvalidRedirectUri(t *testing.T) {
	codeRequest := &AuthorizeTestCase{
		ResponseType: string(RESPONSE_TYPE_CODE),
		ClientId:     CLIENT_ID_CODE_WITH_PKCE,
		RedirectUri:  "http://fail/back",
	}
	testInvalidAuthorizationRequest(t, "Code request with invalid response_type", codeRequest, DESC_INVALID_REDIRECT_URI)
}

/**************************************************/
/*                                                */
/*            IMPLICIT FLOW TEST CASES            */
/*                                                */
/**************************************************/
func TestValidImplicitFlowWithNoState(t *testing.T) {
	//Build a valid implicit call with no state
	request := &AuthorizeTestCase{
		ResponseType: string(RESPONSE_TYPE_TOKEN),
		ClientId:     CLIENT_ID_IMPLICIT,
		RedirectUri:  VALID_REDIRECT_URI,
	}

	testValidImplicitFlow(t, "Valid implicit flow", request)
}

func TestValidImplicitFlowWithState(t *testing.T) {
	//Build a valid implicit url providing a state that would be returned
	request := &AuthorizeTestCase{
		ResponseType: string(RESPONSE_TYPE_TOKEN),
		ClientId:     CLIENT_ID_IMPLICIT,
		RedirectUri:  VALID_REDIRECT_URI,
		State:        "state",
	}

	testValidImplicitFlow(t, "Valid implicit flow", request)
}

func TestImplicitFlowWithoutRedirectUri(t *testing.T) {
	//Build a valid implicit url but without redirect_uri
	request := &AuthorizeTestCase{
		ResponseType: string(RESPONSE_TYPE_TOKEN),
		ClientId:     CLIENT_ID_IMPLICIT,
	}

	testValidImplicitFlow(t, "Valid implicit flow without redirect uri provided into the url", request)
}

func TestInvalidImplicitFlow(t *testing.T) {
	//Build an implicit request with an invalid redirect_uri
	request := &AuthorizeTestCase{
		ResponseType: string(RESPONSE_TYPE_TOKEN),
		ClientId:     CLIENT_ID_IMPLICIT,
		RedirectUri:  "http://fail/back",
	}

	testInvalidAuthorizationRequest(t, "Invalid implicit flow", request, DESC_INVALID_REDIRECT_URI)

}

//common method for valid authorization code flow
func testValidAuthorizationCodeFlow(t *testing.T, testName string, codeRequest *AuthorizeTestCase) {

	rr := callAuthorizationRequestHandler(t, codeRequest)

	//Check the status code is what we expect.
	if http.StatusFound != rr.Code {
		t.Error(unexpectedCodeStatusMessage(testName, rr.Code, http.StatusFound))
	}

	location := rr.HeaderMap.Get("Location")
	uri, _ := url.Parse(location)
	query := uri.Query()
	code := query.Get(PARAM_CODE)

	if code == "" {
		t.Errorf("%v test : Did not find expected (%v) query parameter into this url : %v",
			testName, PARAM_CODE, location)
	}
}

//common method for valid implicit flow
func testValidImplicitFlow(t *testing.T, testName string, implcitRequest *AuthorizeTestCase) {

	rr := callAuthorizationRequestHandler(t, implcitRequest)

	//Check the status code is what we expect.
	if http.StatusFound != rr.Code {
		t.Error(unexpectedCodeStatusMessage(testName, rr.Code, http.StatusFound))
	}

	location := rr.HeaderMap.Get("Location")
	uri, _ := url.Parse(location)

	if errorMessage := "%v test was expecting a fragment %v into the url : %v"; uri.Fragment == "" {
		t.Fatalf(errorMessage, testName, "", location)
	} else {
		//Extract all params from the fragment
		fragments := make(map[string]string)
		for _, fragment := range strings.Split(uri.Fragment, "&") {
			param := strings.Split(fragment, "=")
			fragments[param[0]] = param[1]
		}
		if fragments[PARAM_ACCESS_TOKEN] == "" {
			t.Errorf(errorMessage, testName, PARAM_ACCESS_TOKEN, uri.Fragment)
		}
		if fragments[PARAM_TOKEN_TYPE] == "" {
			t.Errorf(errorMessage, testName, PARAM_TOKEN_TYPE, uri.Fragment)
		}
		if fragments[PARAM_STATE] != implcitRequest.State {
			t.Errorf(errorMessage, testName, PARAM_TOKEN_TYPE, uri.Fragment)
		}
	}
}

func testInvalidAuthorizationRequest(t *testing.T, testName string, request *AuthorizeTestCase, expectedError string) {

	rr := callAuthorizationRequestHandler(t, request)

	if http.StatusBadRequest != rr.Code {
		t.Error(unexpectedCodeStatusMessage(testName, rr.Code, http.StatusBadRequest))
	}

	//If status is not ok, and error message is not as expected
	if body := rr.Body.String(); !strings.Contains(body, expectedError) {
		t.Errorf("%v test : Was expecting error message (%v) but got : %v", testName, expectedError, body)
	}

}

func callAuthorizationRequestHandler(t *testing.T, request *AuthorizeTestCase) *httptest.ResponseRecorder {

	if req, err := http.NewRequest(http.MethodGet, buildTestCaseUrl(request), nil); err != nil {
		t.Fatal(err)
		return nil
	} else {
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		return rr
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
	query.Add(PARAM_STATE, testCase.State)

	uri.RawQuery = query.Encode()
	return uri.String()
}

func unexpectedCodeStatusMessage(testName string, got, expected int) string {
	return fmt.Sprintf("%v test : Unexpected code status, got (%v) while expecting (%v)", testName, got, expected)
}
