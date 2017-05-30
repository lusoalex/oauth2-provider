package oauth2Provider

import (
	"net/http"
	"html/template"
)

const (
	TEMPLATE_ROOT = "templates/"
	ERROR_TEMPLATE = TEMPLATE_ROOT+"error.html"
)

type Error struct {
	Status int
	Message template.HTML
}

var templates = template.Must(template.ParseFiles(ERROR_TEMPLATE))

func handleError(w http.ResponseWriter, err error, status int) {

	errorMessage := &Error{Status: status, Message: template.HTML(err.Error())}

	err = templates.ExecuteTemplate(w, ERROR_TEMPLATE, errorMessage)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
