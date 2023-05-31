package utils

// ErrorCode implements error
var _ error = (*ErrorCode)(nil)

type ErrorCode struct {
	error string
	code  int
}

func NewErrorCode(statusCode int, err error) *ErrorCode {
	return NewErrorCodeString(statusCode, err.Error())
}

func NewErrorCodeString(statusCode int, errMessage string) *ErrorCode {
	return &ErrorCode{
		error: errMessage,
		code:  statusCode,
	}
}

func (e ErrorCode) Error() string {
	return e.error
}

func (e ErrorCode) Code() int {
	return e.code
}
