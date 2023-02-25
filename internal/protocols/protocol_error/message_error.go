package protocol_error

const ErrorMessage ErrorCode = 1

func MessageError(message string) *Error {
	return &Error{
		Code:    ErrorMessage,
		Message: message,
	}
}
