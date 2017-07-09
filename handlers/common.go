package handlers

import (
	"errors"
	"log"
	"net/http"
	"oauth2-provider/constants"
	"oauth2-provider/response"
	"path"
	"strings"
	"time"
)

var NotFound = errors.New("Not found")

func ShiftPath(p string) (head, tail string) {
	p = path.Clean(p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

type WrappedResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *WrappedResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

type CommonHandler struct{}

func (h *CommonHandler) ServeHTTP(_w http.ResponseWriter, req *http.Request) {
	w := &WrappedResponseWriter{
		ResponseWriter: _w,
		StatusCode:     http.StatusOK,
	}
	defer func(path string, start time.Time) {
		log.Printf("%s %s %d %s", req.Method, path, w.StatusCode, time.Since(start))
	}(req.URL.String(), time.Now())

	var response response.Response

	if response, err := h.Handle(w, req); err != nil {
		switch err {
		case NotFound:
			http.NotFound(w, req)
		default:
			http.Error(w, "Internal  server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set(constants.CONTENT_TYPE, response.ContentType())
	bytes, _ := response.Render() //TODO manage error
	w.Write(bytes)
}

func (h *CommonHandler) Handle(w http.ResponseWriter, req *http.Request) (response.Response, error) {
	return nil, NotFound
}
