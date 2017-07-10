package response

import (
	"encoding/json"
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
	return &Json{Content: content}
}
