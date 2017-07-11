package models

type Error struct {
	Reason           string `json:"error"`                       //required
	ErrorDescription string `json:"error_description,omitempty"` //Optional
	ErrorUri         string `json:"error_uri,omitempty"`         //Optional
	State            string `json:"state,omitempty"`             //Required if present into the request.
}

type BadRequest struct {
	Error
}

type ForbiddenRequest struct {
	Error
}

func (e *Error) Error() string {
	return e.Reason + " : " + e.ErrorDescription
}
