package oauth2Provider

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOauth2Handler(t *testing.T) {

	//Test all declared endpoints.
	testCases := []DataTestCase{
		{name: "Health_check endpoint", url: "/health_check", status: http.StatusOK},
		{name: "Authorize endpoint", method: "GET", url: "/authorize", status: http.StatusBadRequest},
		{name: "Token endpoint", method: "POST", url: "/token", status: http.StatusBadRequest},
		{name: "Not Found test", method: "GET", url: "/notfound", status: http.StatusNotFound},
		{name: "Method Not Allowed test", method: "GET", url: "/token", status: http.StatusMethodNotAllowed},
	}

	server := httptest.NewServer(Oauth2Handler())
	defer server.Close() //close when test is ended...

	for _, test := range testCases {

		var req *http.Request
		var err error

		if req, err = http.NewRequest(test.method, server.URL+test.url, nil); err != nil {
			t.Fatal(err)
		}

		if resp, err := http.DefaultClient.Do(req); err != nil {
			t.Fatal(err)
		} else {
			//Check the status code is what we expect.
			assertStatus(t, &test, resp.StatusCode)
		}
	}
}
