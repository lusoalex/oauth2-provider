package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"oauth2-provider/models"
	"path"
	"strings"
)

/*
 * Oauth2 request may return different error response according to : https://tools.ietf.org/html/rfc6749#section-5.2
 */
type oauth2RequestHandler func(http.ResponseWriter, *http.Request) error

/*
 * Manage oauth error response
 */
func serveOauth2Request(w http.ResponseWriter, r *http.Request, handlerWithError oauth2RequestHandler) {

	err := handlerWithError(w, r)

	if err != nil {

		if body, marshalErr := json.Marshal(err); marshalErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		} else {
			switch err.(type) {
			case *models.BadRequest:
				w.WriteHeader(http.StatusBadRequest)
			case *models.ForbiddenRequest:
				w.WriteHeader(http.StatusForbidden)
			default:
				log.Printf("Non Oauth2 error type : [%T], value[%v]\n", err, err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Add("Content-type", "application/json")
			w.Write(body)
		}
	}
}

func ShiftPath(p string) (head, tail string) {
	p = path.Clean(p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
