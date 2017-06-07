package oauth2Provider

import (
	"fmt"
	"net/http"
	"regexp"
)

type MyOauth2Handler struct{}

var validPath = regexp.MustCompile("^/([a-zA-Z0-9_]+)?.*")

func (h *MyOauth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	uri := r.URL.Path

	m := validPath.FindStringSubmatch(uri)
	if m == nil {
		fmt.Println("aie aie aie aie")
		return
	}

	switch m[1] {
	case "health_check":
		handleHealthCheck(w, r)
	case "authorize":
		handleAuthorizationRequest(w, r)
	//case "token" : handleTokenRequest(w,r)
	default:
		http.Error(w, "No matching resource found", http.StatusNotFound)
	}

	return
}

func LaunchServer() {
	handler := new(MyOauth2Handler)
	http.ListenAndServe(":8000", handler)
}
