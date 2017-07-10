package response

import (
	"net/http"
)

type HTTPResponse struct {
	Status int
	ContentType string
	Body []byte
}

func NewHTTPResponse(body []byte, contentType string) *HTTPResponse {
	return &HTTPResponse{
		Body: body,
		Status: http.StatusOK,
		ContentType: contentType,
	}
}

type Response interface {
	Render() (*HTTPResponse, error)
}

// func (r *ResponseStatus) Status() int {
// 	if r.httpStatus == 0 {
// 		return http.StatusOK
// 	}
// 	return r.httpStatus
// }

// func (r *ResponseStatus) OK() *Response {
// 	r.httpStatus = http.StatusOK
// 	return r
// }

// func (r *ResponseStatus) BadRequest() *Response {
// 	r.httpStatus = http.StatusBadRequest
// 	return r
// }

// func (r *ResponseStatus) NotImplemented() *Response {
// 	r.httpStatus = http.StatusNotImplemented
// 	return r
// }

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

