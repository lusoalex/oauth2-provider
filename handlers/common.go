package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"oauth2-provider/models"
	"oauth2-provider/utils"
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

func inRange(key string, list []string) bool {
	for _, val := range list {
		if val == key {
			return true
		}
	}
	return false
}

type Oauth2Handler struct {
	utils.KeyValueStore
}

type oauth2RequestHandler func(http.ResponseWriter, *http.Request) error

func (*Oauth2Handler) handleOauth2Request(path string, allowedMethods []string, w http.ResponseWriter, r *http.Request, handler oauth2RequestHandler) {

	if head, _ := ShiftPath(r.URL.Path); head != path || !inRange(r.Method, allowedMethods) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

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
