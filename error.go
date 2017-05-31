package oauth2Provider

import (
	"html/template"
	"net/http"
)

const (
	TEMPLATE_PATH  = "templates/"
	TEMPLATE_ERROR = "error.html"
)

type Error struct {
	Status  int
	Message template.HTML
}

var templates = template.Must(template.ParseFiles(TEMPLATE_PATH + TEMPLATE_ERROR))

func handleError(w http.ResponseWriter, err error, status int) {

	errorMessage := &Error{Status: status, Message: template.HTML(err.Error())}

	w.WriteHeader(status)

	err = templates.ExecuteTemplate(w, TEMPLATE_ERROR, errorMessage)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
