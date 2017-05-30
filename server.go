package oauth2Provider

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

type MyOauth2Handler struct{}

var validPath = regexp.MustCompile("^/([a-zA-Z0-9_]+)?.*")

func (h *MyOauth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	uri := r.URL.Path
	fmt.Println(uri)

	m := validPath.FindStringSubmatch(uri)
	if m == nil {
		fmt.Println("aie aie aie aie")
		return
	}

	fmt.Println(m[1])

	switch m[1] {
	case "health_check":
		handleHealthCheck(w, r)
	case "authorize":
		handleAuthorizationRequest(w, r)
	//case "token" : handleTokenRequest(w,r)
	default:
		handleError(w, errors.New("No matching resource found"), http.StatusNotFound)
	}

	fmt.Printf("%v",w)

	return
}

func LaunchServer() {
	handler := new(MyOauth2Handler)
	http.ListenAndServe(":8000", handler)
}

/*
func (h *MyHandlerType) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	// ...use uri...
}
*/

/*

import (
	"html/template"
	"net/http"
	"io/ioutil"
	"regexp"
)

var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, "templates/"+tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

*/
