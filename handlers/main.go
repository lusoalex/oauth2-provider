package handlers

import (
	"log"
	"net/http"
	"time"
)

type MainHandler struct{}

func (h *MainHandler) ServeHTTP(_w http.ResponseWriter, req *http.Request) {

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
		(&HealthCheckHandler{}).ServeHTTP(w, req)
	case "authorize":
		(&AuthorizeHandler{}).ServeHttp(w, req)

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
