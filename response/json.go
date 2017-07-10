package response

import (
	"net/http"
	"encoding/json"
	"oauth2-provider/constants"
)

type Json struct {
	Content interface{}
}

func (r *Json) Render() (*HTTPResponse, error) {
	body, err := json.Marshal(r.Content)
	if err != nil {
		return nil, err
	}
	return NewHTTPResponse(body, "application/json"), nil
}

func NewJsonResponse(content interface{}) *Json {
	return &JsonResponse{Content: content}
}
