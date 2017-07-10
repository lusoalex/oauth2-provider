package response

import (
	"encoding/json"
	"net/http"
	"oauth2-provider/constants"
)

type JsonResponse struct {
	Content    interface{}
	httpStatus int
}

func (r *JsonResponse) Render() ([]byte, error) {
	return json.Marshal(r.Content)
}

func (r *JsonResponse) ContentType() string {
	return constants.CONTENT_TYPE_JSON
}

func (r *JsonResponse) Status() int {
	return r.httpStatus
}

func NewJsonResponse(content interface{}) *JsonResponse {
	return &JsonResponse{Content: content, httpStatus: http.StatusOK}
}

func (r *JsonResponse) BadRequest() *JsonResponse {
	r.httpStatus = http.StatusBadRequest
	return r
}
