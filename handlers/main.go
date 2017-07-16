package handlers

import (
	"log"
	"net/http"
	"oauth2-provider/settings"
	"time"
)

type MainHandler struct {
	Opts *settings.Oauth2ProviderSettings
}

func (main *MainHandler) ServeHTTP(_w http.ResponseWriter, req *http.Request) {

	w := &WrappedResponseWriter{
		ResponseWriter: _w,
		StatusCode:     http.StatusOK,
	}
	defer func(path string, start time.Time) {
		log.Printf("%s %s %d %s", req.Method, path, w.StatusCode, time.Since(start))
	}(req.URL.String(), time.Now())

	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	switch head {
	case "health_check":
		(&HealthCheckHandler{main.Opts}).ServeHTTP(w, req)
	case "authorize":
		(&AuthorizeHandler{main.Opts}).ServeHTTP(w, req)
	case "token":
		(&TokenHandler{main.Opts}).ServeHTTP(w, req)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
}

type WrappedResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *WrappedResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}
