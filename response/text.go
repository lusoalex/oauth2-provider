package response

type Text struct {
	Content string
}

func (r *Text) Render() (*HTTPResponse, error) {
	return NewHTTPResponse([]byte(r.Content), "text/plain"), nil
}

