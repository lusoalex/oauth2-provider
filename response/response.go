package response

import (
	"net/http"
)

type HTTPResponse struct {
	Status int
	ContentType string
	Body []byte
	Header map[string]string
}

func NewHTTPResponse(body []byte, contentType string) *HTTPResponse {
	return &HTTPResponse{
		Body: body,
		Status: http.StatusOK,
		ContentType: contentType,
		Header:make(map[string]string),
	}
}

type Response interface {
	Render() (*HTTPResponse, error)
}

func OK(response Response) (*HTTPResponse, error) {
	r, err := response.Render()
	if err != nil {
		return nil, err
	} else {
		r.Status = http.StatusOK
		return r, nil
	}
}

func BadRequest(response Response) (*HTTPResponse, error) {
	r, err := response.Render()
	if err != nil {
		return nil, err
	} else {
		r.Status = http.StatusBadRequest
		return r, nil
	}
}

func (r *HTTPResponse) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", r.ContentType)
	w.WriteHeader(r.Status)
	w.Write(r.Body)
}

