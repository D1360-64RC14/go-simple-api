package dtos

type ErrorMessage struct {
	Error string `json:"error"`
}

func NewErrorMessage(err error) *ErrorMessage {
	return &ErrorMessage{err.Error()}
}

func NewErrorMessageByString(errMsg string) *ErrorMessage {
	return &ErrorMessage{errMsg}
}
