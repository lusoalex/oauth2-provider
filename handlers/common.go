package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"oauth2-provider/models"
	"path"
	"strings"
)

func ShiftPath(p string) (head, tail string) {
	p = path.Clean(p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

type oauth2RequestHandler func(http.ResponseWriter, *http.Request) error

func HandleOauth2Request(w http.ResponseWriter, r *http.Request, handler oauth2RequestHandler) {
	err := handler(w, r)

	if err != nil {

		log.Printf("Oauth2Error type [%T], value[%v]\n", err, err)

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
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Add("Content-type", "application/json")
			w.Write(body)
		}
	}
}
