package response

import (
	"encoding/json"
	"oauth2-provider/constants"
)

type JsonResponse struct {
	Content interface{}
}

func (r *JsonResponse) Render() ([]byte, error) {
	return json.Marshal(r.Content)
}

func (r *JsonResponse) ContentType() string {
	return constants.CONTENT_TYPE_JSON
}
