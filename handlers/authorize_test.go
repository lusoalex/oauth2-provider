package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/*
func (a *AuthorizeHandlerTest) handleAuthorizationRequest(w http.ResponseWriter, r *http.Request) error {
	return nil
}
*/

func TestServeHttp(t *testing.T) {

	testCase := []struct {
		method string
		path   string
		code   int
		body   string
	}{
		{method: http.MethodGet, path: "/unknown", code: 404, body: "Not found"},
		{method: http.MethodConnect, path: "/authorize", code: 404, body: "Not found"},
		{method: http.MethodDelete, path: "/authorize", code: 404, body: "Not found"},
		{method: http.MethodHead, path: "/authorize", code: 404, body: "Not found"},
		{method: http.MethodOptions, path: "/authorize", code: 404, body: "Not found"},
		{method: http.MethodPatch, path: "/authorize", code: 404, body: "Not found"},
		{method: http.MethodPut, path: "/authorize", code: 404, body: "Not found"},
		{method: http.MethodTrace, path: "/authorize", code: 404, body: "Not found"},
		{method: http.MethodGet, path: "/authorize", code: 400, body: "error_description"},
		{method: http.MethodPost, path: "/authorize", code: 400, body: "error_description"},
	}

	handler := http.HandlerFunc((&AuthorizeHandler{}).ServeHTTP)

	for _, test := range testCase {
		if req, err := http.NewRequest(test.method, test.path, nil); err != nil {
			t.Fatal(err)
		} else {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if test.code != rr.Code {
				t.Errorf("Unexpected response status code, got %v while expecting %v\n", rr.Code, test.code)
			}
			if !strings.Contains(rr.Body.String(), test.body) {
				t.Errorf("Unexpected response body, was expecting at least this value %v ", test.body)
			}
		}
	}
}

func TestHandleAuthorizationRequest(t *testing.T) {

	rw := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/authorize?client_id=toto", nil)
	err := (&AuthorizeHandler{}).handleAuthorizationRequest(rw, r)
	fmt.Println(err)
	fmt.Println(rw.Code)
	fmt.Println(rw.Body)
}
