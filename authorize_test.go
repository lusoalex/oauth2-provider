package oauth2Provider

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestHandleAuthorizationRequest(t *testing.T) {

	//Initialize request
	req,err := http.NewRequest("GET", "http://localhost:8000/authorize", nil)

	if err != nil {
		t.Fatal(err)
	}

	//Add query parameters
	req.Form.Set(PARAM_RESPONSE_TYPE,string(CODE))
	req.Form.Set(PARAM_CLIENT_ID,"toto")
	req.Form.Set(PARAM_REDIRECT_URI,"http://callback")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleHealthCheck)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}


}