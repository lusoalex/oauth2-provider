package response

import "encoding/json"

type JsonResponse struct {
	Content interface{}
}

func (r *JsonResponse) Render() ([]byte, error) {
	return json.Marshal(r.Content)
}
