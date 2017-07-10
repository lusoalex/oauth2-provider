package response

import "net/http"

type Response interface {
	Render() ([]byte, error)
	ContentType() string
	Status() int
}

type ResponseStatus struct {
	httpStatus int
}

func (r *ResponseStatus) Status() int {
	if r.httpStatus == 0 {
		return http.StatusOK
	}
	return r.httpStatus
}

func (r *ResponseStatus) OK() *Response {
	r.httpStatus = http.StatusOK
	return r
}

func (r *ResponseStatus) BadRequest() *Response {
	r.httpStatus = http.StatusBadRequest
	return r
}
