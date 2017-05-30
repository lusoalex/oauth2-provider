package oauth2Provider

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestServeHTTP(t *testing.T) {

	handler := new(MyOauth2Handler)
	server := httptest.NewServer(handler)
	defer server.Close()

	resp,err := http.Get("http://localhost:8000/health_check")

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("Expecting Content-Type value : application/json, but was : %v",contentType)
	}
}