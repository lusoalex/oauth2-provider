package oauth2Provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTP(t *testing.T) {

	handler := new(MyOauth2Handler)
	server := httptest.NewServer(handler)
	defer server.Close() //close when test is ended...

	resp, err := http.Get(server.URL + "/health_check")

	fmt.Println("server url is ", server.URL)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	if contentType := resp.Header.Get(CONTENT_TYPE); contentType != CONTENT_TYPE_JSON {
		t.Fatalf("Expecting Content-Type value : %v, but was : %v", CONTENT_TYPE_JSON, contentType)
	}
}
